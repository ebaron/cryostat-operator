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
	"fmt"
	"time"

	operatorv1beta1 "github.com/cryostatio/cryostat-operator/api/v1beta1"
	scapiv1alpha3 "github.com/operator-framework/api/pkg/apis/scorecard/v1alpha3"
	apimanifests "github.com/operator-framework/api/pkg/manifests"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	OperatorInstallTestName string        = "operator-install"
	CryostatCRTestName      string        = "cryostat-cr"
	operatorDeploymentName  string        = "cryostat-operator-controller-manager"
	testTimeout             time.Duration = time.Minute * 10
)

// OperatorInstallTest checks that the operator installed correctly
func OperatorInstallTest(bundle *apimanifests.Bundle, namespace string) scapiv1alpha3.TestResult {
	r := scapiv1alpha3.TestResult{}
	r.Name = OperatorInstallTestName
	r.State = scapiv1alpha3.PassState
	r.Errors = make([]string, 0)
	r.Suggestions = make([]string, 0)

	// Create a new Kubernetes REST client for this test
	client, err := NewClientset()
	if err != nil {
		return fail(r, fmt.Sprintf("failed to create client: %s", err.Error()))
	}

	// Poll the deployment until it becomes available or we timeout
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	err = waitForDeploymentAvailability(ctx, client, namespace, operatorDeploymentName, &r)
	if err != nil {
		return fail(r, fmt.Sprintf("operator deployment did not become available: %s", err.Error()))
	}

	return r
}

// CryostatCRTest checks that the operator installs Cryostat in response to a Cryostat CR
func CryostatCRTest(bundle *apimanifests.Bundle, namespace string) scapiv1alpha3.TestResult {
	r := scapiv1alpha3.TestResult{}
	r.Name = CryostatCRTestName
	r.State = scapiv1alpha3.PassState
	r.Errors = make([]string, 0)
	r.Suggestions = make([]string, 0)

	// Create a new Kubernetes REST client for this test
	client, err := NewClientset()
	if err != nil {
		return fail(r, fmt.Sprintf("failed to create client: %s", err.Error()))
	}

	// Create a default Cryostat CR
	cr := &operatorv1beta1.Cryostat{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cryostat-cr-test",
			Namespace: namespace,
		},
		Spec: operatorv1beta1.CryostatSpec{
			Minimal: false,
		},
	}

	ctx := context.Background()
	cr, err = client.OperatorCRDs().Cryostats(namespace).Create(ctx, cr)
	if err != nil {
		return fail(r, fmt.Sprintf("failed to create Cryostat CR: %s", err.Error()))
	}

	// Poll the deployment until it becomes available or we timeout
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	err = waitForDeploymentAvailability(ctx, client, cr.Namespace, cr.Name, &r)
	if err != nil {
		return fail(r, fmt.Sprintf("Cryostat main deployment did not become available: %s", err.Error()))
	}

	err = wait.PollImmediateUntilWithContext(ctx, time.Second, func(ctx context.Context) (done bool, err error) {
		cr, err = client.OperatorCRDs().Cryostats(namespace).Get(ctx, cr.Name)
		if err != nil {
			return false, fmt.Errorf("failed to get Cryostat CR: %s", err.Error())
		}
		if len(cr.Status.ApplicationURL) > 0 {
			return true, nil
		}
		r.Log += "Application URL is not yet available\n"
		return false, nil
	})
	if err != nil {
		return fail(r, fmt.Sprintf("Application URL not found in CR: %s", err.Error()))
	}
	r.Log += fmt.Sprintf("Application is ready at %s\n", cr.Status.ApplicationURL)

	return r
}

func waitForDeploymentAvailability(ctx context.Context, client *CryostatClientset, namespace string,
	name string, r *scapiv1alpha3.TestResult) error {
	return wait.PollImmediateUntilWithContext(ctx, time.Second, func(ctx context.Context) (done bool, err error) {
		deploy, err := client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			if kerrors.IsNotFound(err) {
				r.Log += fmt.Sprintf("deployment %s is not yet found\n", name)
				return false, nil // Retry
			}
			return false, fmt.Errorf("failed to get deployment: %s", err.Error())
		}
		// Check for Available condition
		for _, condition := range deploy.Status.Conditions {
			if condition.Type == appsv1.DeploymentAvailable &&
				condition.Status == corev1.ConditionTrue {
				r.Log += fmt.Sprintf("deployment %s is available\n", deploy.Name)
				return true, nil
			}
		}
		r.Log += fmt.Sprintf("deployment %s is not yet available\n", deploy.Name)
		return false, nil
	})
}

func fail(r scapiv1alpha3.TestResult, message string) scapiv1alpha3.TestResult {
	r.State = scapiv1alpha3.FailState
	r.Errors = append(r.Errors, message)
	return r
}
