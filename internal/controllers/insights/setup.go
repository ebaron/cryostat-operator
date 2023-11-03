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

package insights

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/cryostatio/cryostat-operator/internal/controllers/common"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type InsightsIntegration struct {
	Manager ctrl.Manager
	Log     *logr.Logger
	common.OSUtils
}

func NewInsightsIntegration(mgr ctrl.Manager, log *logr.Logger) *InsightsIntegration {
	return &InsightsIntegration{
		Manager: mgr,
		Log:     log,
		OSUtils: &common.DefaultOSUtils{},
	}
}

func (i *InsightsIntegration) Setup() (*url.URL, error) {
	var proxyUrl *url.URL
	namespace := i.getOperatorNamespace()
	// This will happen when running the operator locally
	if len(namespace) == 0 {
		i.Log.Info("Operator namespace not detected, disabling Insights integration")
		return nil, nil
	}

	ctx := context.Background()
	if i.isInsightsEnabled() {
		err := i.createInsightsController(namespace)
		if err != nil {
			i.Log.Error(err, "unable to add controller to manager", "controller", "Insights")
			return nil, err
		}
		// Create a Config Map to be used as a parent of all Insights Proxy related objects
		err = i.createConfigMap(ctx, namespace)
		if err != nil {
			i.Log.Error(err, "failed to create config map for Insights")
			return nil, err
		}
		proxyUrl = i.getProxyURL(namespace)
	} else {
		// Delete any previously created Config Map (and its children)
		err := i.deleteConfigMap(ctx, namespace)
		if err != nil {
			i.Log.Error(err, "failed to delete config map for Insights")
			return nil, err
		}

	}
	return proxyUrl, nil
}

func (i *InsightsIntegration) isInsightsEnabled() bool {
	return strings.ToLower(i.GetEnv(EnvInsightsEnabled)) == "true"
}

func (i *InsightsIntegration) getOperatorNamespace() string {
	return i.GetEnv("NAMESPACE")
}

func (i *InsightsIntegration) createInsightsController(namespace string) error {
	config := &InsightsReconcilerConfig{
		Client:    i.Manager.GetClient(),
		Log:       ctrl.Log.WithName("controllers").WithName("Insights"),
		Scheme:    i.Manager.GetScheme(),
		Namespace: namespace,
		OSUtils:   i.OSUtils,
	}
	controller, err := NewInsightsReconciler(config)
	if err != nil {
		return err
	}
	if err := controller.SetupWithManager(i.Manager); err != nil {
		return err
	}
	return nil
}

func (i *InsightsIntegration) createConfigMap(ctx context.Context, namespace string) error {
	// The config map should be owned by the operator deployment to ensure it and its descendants are garbage collected
	owner := &appsv1.Deployment{}
	err := i.Manager.GetAPIReader().Get(ctx, types.NamespacedName{Name: OperatorDeploymentName, Namespace: namespace}, owner)
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InsightsConfigMapName,
			Namespace: namespace,
		},
	}
	err = controllerutil.SetControllerReference(owner, cm, i.Manager.GetScheme())
	if err != nil {
		return err
	}

	err = i.Manager.GetClient().Create(ctx, cm, &client.CreateOptions{})
	if err == nil {
		i.Log.Info("Config Map for Insights created", "name", cm.Name, "namespace", cm.Namespace)
	}
	// This may already exist if the pod restarted
	return client.IgnoreAlreadyExists(err)
}

func (i *InsightsIntegration) deleteConfigMap(ctx context.Context, namespace string) error {
	// Children will be garbage collected
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InsightsConfigMapName,
			Namespace: namespace,
		},
	}

	err := i.Manager.GetClient().Delete(ctx, cm, &client.DeleteOptions{})
	if err == nil {
		i.Log.Info("Config Map for Insights deleted", "name", cm.Name, "namespace", cm.Namespace)
	}
	// This may not exist if no config map was previously created
	return client.IgnoreNotFound(err)
}

func (i *InsightsIntegration) getProxyURL(namespace string) *url.URL {
	return &url.URL{
		Scheme: "http", // TODO add https support (r.IsCertManagerInstalled)
		Host:   fmt.Sprintf("%s.%s.svc.cluster.local", ProxyServiceName, namespace),
	}
}
