// Copyright The Cryostat Authors
//
// The Universal Permissive License (UPL), Version 1.0
//
// Subject to the condition set forth below, permission is hereby granted to any
// person obtaining a copy of this software, associated documentation and/or data
// (collectively the "Software"), free of charge and under any and all copyright
// rights in the Software, and any and all patent rights owned or freely
// licensable by each licensor hereunder covering either (i) the unmodified
// Software as contributed to or provided by such licensor, or (ii) the Larger
// Works (as defined below), to deal in both
//
// (a) the Software, and
// (b) any piece of software and/or hardware listed in the lrgrwrks.txt file if
// one is included with the Software (each a "Larger Work" to which the Software
// is contributed by such licensors),
//
// without restriction, including without limitation the rights to copy, create
// derivative works of, display, perform, and distribute the Software and make,
// use, sell, offer for sale, import, export, have made, and have sold the
// Software and the Larger Work(s), and to sublicense the foregoing rights on
// either these or other terms.
//
// This license is subject to the following condition:
// The above copyright notice and either this complete permission notice or at
// a minimum a reference to the UPL must be included in all copies or
// substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package scorecard

import (
	"context"
	"errors"
	"fmt"

	"github.com/blang/semver/v4"
	scapiv1alpha3 "github.com/operator-framework/api/pkg/apis/scorecard/v1alpha3"
	operatorsv1 "github.com/operator-framework/api/pkg/operators/v1"
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"

	configv1 "github.com/openshift/api/config/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
)

func installOpenShiftCertManager(r *scapiv1alpha3.TestResult) error {
	ctx := context.Background()

	// Get in-cluster REST config from pod
	config, err := ctrl.GetConfig()
	if err != nil {
		return err
	}
	scheme := runtime.NewScheme()
	if err := operatorsv1alpha1.AddToScheme(scheme); err != nil {
		return err
	}
	if err := operatorsv1.AddToScheme(scheme); err != nil {
		return err
	}
	if err := configv1.AddToScheme(scheme); err != nil {
		return err
	}

	// Client to use with ClusterVersion and OperatorHub API
	openshiftClient, err := newRESTClientForGV(config, scheme, &configv1.SchemeGroupVersion)
	if err != nil {
		return err
	}

	// Check which OpenShift version we're on. Different versions of the cert-manager operator
	// support different install modes. We have to differentiate when we create the subscription.
	clusterVersions := &configv1.ClusterVersionList{}
	err = openshiftClient.Get().Resource("clusterversions").Do(ctx).Into(clusterVersions)
	if err != nil {
		return err
	}
	if len(clusterVersions.Items) == 0 {
		return errors.New("no ClusterVersions found")
	}
	clusterVersion := clusterVersions.Items[0]
	r.Log += fmt.Sprintf("OpenShift version is: %s\n", clusterVersion.Status.Desired.Version)
	trimmedVer, err := semver.FinalizeVersion(clusterVersion.Status.Desired.Version)
	if err != nil {
		return err
	}
	version, err := semver.Parse(trimmedVer)
	if err != nil {
		return err
	}
	// The stable-v1 channel is available on OpenShift 4.12+
	useStable := false
	if version.GTE(semver.MustParse("4.12.0")) {
		useStable = true
	}

	// Patch the OperatorHub config to enable the default catalog sources
	hubPatch := `[{"op": "add", "path": "/spec/disableAllDefaultSources", "value": false}]`
	hub := &configv1.OperatorHub{}
	err = openshiftClient.Patch(types.JSONPatchType).Resource("operatorhubs").Name("cluster").Body([]byte(hubPatch)).Do(ctx).Into(hub)
	if err != nil {
		return err
	}
	r.Log += "OperatorHub patched to enable default catalog sources\n"

	// Client for Subscription and ClusterServiceVersion APIs
	olmClient, err := newRESTClientForGV(config, scheme, &operatorsv1alpha1.SchemeGroupVersion)
	if err != nil {
		return err
	}

	// With the stable-v1 channel, we need to install the operator into one
	// namespace. This requires us to create the namespace and an OperatorGroup for it.
	subNamespace := "openshift-operators"
	channel := "tech-preview"
	if useStable {
		subNamespace = "cert-manager-operator"
		channel = "stable-v1"

		// Client for Namespaces
		client, err := corev1client.NewForConfig(config)
		if err != nil {
			return err
		}
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "cert-manager-operator",
			},
		}
		_, err = client.Namespaces().Create(ctx, ns, metav1.CreateOptions{})
		if err != nil && !kerrors.IsAlreadyExists(err) {
			return err
		}
		r.Log += fmt.Sprintf("Created %s namespace\n", subNamespace)

		og := &operatorsv1.OperatorGroup{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cert-manager-og",
				Namespace: subNamespace,
			},
			Spec: operatorsv1.OperatorGroupSpec{
				TargetNamespaces: []string{
					subNamespace,
				},
			},
		}

		// Client for OperatorGroup API
		ogClient, err := newRESTClientForGV(config, scheme, &operatorsv1.SchemeGroupVersion)
		if err != nil {
			return err
		}
		err = ogClient.Post().Resource("operatorgroups").Namespace(og.Namespace).Name(og.Name).Body(og).Do(ctx).Into(og)
		if err != nil && !kerrors.IsAlreadyExists(err) {
			return err
		}
		r.Log += fmt.Sprintf("Created OperatorGroup for %s\n", subNamespace)
	}

	// Create the Subscription for the cert-manager operator. The namespace and channel
	// are dependent on the OpenShift version
	sub := &operatorsv1alpha1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cert-manager-sub",
			Namespace: subNamespace,
		},
		Spec: &operatorsv1alpha1.SubscriptionSpec{
			CatalogSource:          "redhat-operators",
			CatalogSourceNamespace: "openshift-marketplace",
			Package:                "openshift-cert-manager-operator",
			Channel:                channel,
		},
	}
	err = olmClient.Post().Resource("subscriptions").Namespace(sub.Namespace).Name(sub.Name).Body(sub).Do(ctx).Into(&operatorsv1alpha1.Subscription{})
	if err != nil && !kerrors.IsAlreadyExists(err) {
		return err
	}
	r.Log += fmt.Sprintf("Created Subscription for openshift-cert-manager-operator in %s\n", subNamespace)

	// Check CSV status until we know cert-manager installed successfully
	return wait.PollImmediateUntilWithContext(ctx, testTimeout, func(ctx context.Context) (bool, error) {
		err := olmClient.Get().Resource("subscriptions").Namespace(sub.Namespace).Name(sub.Name).Do(ctx).Into(sub)
		if err != nil {
			return false, fmt.Errorf("failed to get Subscription: %s", err.Error())
		}
		if len(sub.Status.CurrentCSV) == 0 {
			r.Log += fmt.Sprintf("ClusterServiceVersion is not yet found for Subscription %s\n", sub.Name)
			return false, nil
		}

		csv := &operatorsv1alpha1.ClusterServiceVersion{}
		err = olmClient.Get().Resource("clusterserviceversions").Namespace(sub.Namespace).Name(sub.Status.InstalledCSV).Do(ctx).Into(csv)
		if err != nil {
			return false, fmt.Errorf("failed to get ClusterServiceVersion: %s", err.Error())
		}
		// Check for Succeeded phase
		if csv.Status.Phase == operatorsv1alpha1.CSVPhaseSucceeded &&
			csv.Status.Reason == operatorsv1alpha1.CSVReasonInstallSuccessful {
			r.Log += fmt.Sprintf("ClusterServiceVersion %s successfully installed\n", csv.Name)
			return true, nil
		}
		r.Log += fmt.Sprintf("ClusterServiceVersion %s is not yet ready. Current status %s\n", csv.Name, csv.Status.Message)
		return false, nil
	})
}
