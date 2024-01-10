// Copyright The Cryostat Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/url"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	"k8s.io/client-go/discovery"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	certv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	configv1 "github.com/openshift/api/config/v1"
	consolev1 "github.com/openshift/api/console/v1"
	routev1 "github.com/openshift/api/route/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	operatorv1beta1 "github.com/cryostatio/cryostat-operator/api/v1beta1"
	"github.com/cryostatio/cryostat-operator/internal/controllers"
	"github.com/cryostatio/cryostat-operator/internal/controllers/common"
	"github.com/cryostatio/cryostat-operator/internal/controllers/insights"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(operatorv1beta1.AddToScheme(scheme))

	// Register third-party types
	utilruntime.Must(routev1.AddToScheme(scheme))
	utilruntime.Must(certv1.AddToScheme(scheme))
	utilruntime.Must(consolev1.AddToScheme(scheme))
	utilruntime.Must(configv1.AddToScheme(scheme))

	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var enableHTTP2 bool
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&enableHTTP2, "enable-http2", false, "If HTTP/2 should be enabled for the metrics and webhook servers.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	// Customize the webhook server to disable HTTP/2 by default
	disableHTTP2 := func(c *tls.Config) {
		if enableHTTP2 {
			return
		}
		c.NextProtos = []string{"http/1.1"}
	}
	webhookServer := &webhook.Server{
		TLSOpts: []func(config *tls.Config){disableHTTP2},
	}

	// FIXME Disable metrics until this issue is resolved:
	// https://github.com/operator-framework/operator-sdk/issues/4684
	metricsAddr = "0"
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "d696d7ab.redhat.com",
		WebhookServer:          webhookServer,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// To look for particular resource types on the API server
	dc, err := discovery.NewDiscoveryClientForConfig(mgr.GetConfig())
	if err != nil {
		setupLog.Error(err, "failed to create discovery client")
		os.Exit(1)
	}

	openShift, err := isOpenShift(dc)
	if err != nil {
		setupLog.Error(err, "could not determine whether manager is running on OpenShift")
		os.Exit(1)
	}
	environment := "Kubernetes"
	if openShift {
		environment = "OpenShift"
	}
	setupLog.Info(fmt.Sprintf("detected %s environment", environment))

	certManager, err := isCertManagerInstalled(dc)
	if err != nil {
		setupLog.Error(err, "could not determine whether cert-manager is installed")
		os.Exit(1)
	}
	if certManager {
		setupLog.Info("found cert-manager installation")
	} else {
		setupLog.Info("did not find cert-manager installation")
	}

	// Optionally enable Insights integration. Will only be enabled if INSIGHTS_ENABLED is true
	var insightsURL *url.URL
	if openShift {
		insightsURL, err = insights.NewInsightsIntegration(mgr, &setupLog).Setup()
		if err != nil {
			setupLog.Error(err, "failed to set up Insights integration")
		}
	}

	config := newReconcilerConfig(mgr, "ClusterCryostat", "clustercryostat-controller", openShift,
		certManager, insightsURL)
	clusterController, err := controllers.NewClusterCryostatReconciler(config)
	if err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ClusterCryostat")
		os.Exit(1)
	}
	if err = clusterController.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to add controller to manager", "controller", "ClusterCryostat")
		os.Exit(1)
	}
	config = newReconcilerConfig(mgr, "Cryostat", "cryostat-controller", openShift, certManager,
		insightsURL)
	controller, err := controllers.NewCryostatReconciler(config)
	if err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Cryostat")
		os.Exit(1)
	}
	if err = controller.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to add controller to manager", "controller", "Cryostat")
		os.Exit(1)
	}
	if err = (&operatorv1beta1.Cryostat{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "Cryostat")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager", "version", controllers.OperatorVersion)
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func isOpenShift(client discovery.DiscoveryInterface) (bool, error) {
	return discovery.IsResourceEnabled(client, routev1.GroupVersion.WithResource("routes"))
}

func isCertManagerInstalled(client discovery.DiscoveryInterface) (bool, error) {
	return discovery.IsResourceEnabled(client, certv1.SchemeGroupVersion.WithResource("issuers"))
}

func newReconcilerConfig(mgr ctrl.Manager, logName string, eventRecorderName string, openShift bool,
	certManager bool, insightsURL *url.URL) *controllers.ReconcilerConfig {
	return &controllers.ReconcilerConfig{
		Client:                 mgr.GetClient(),
		Log:                    ctrl.Log.WithName("controllers").WithName(logName),
		Scheme:                 mgr.GetScheme(),
		IsOpenShift:            openShift,
		IsCertManagerInstalled: certManager,
		EventRecorder:          mgr.GetEventRecorderFor(eventRecorderName),
		RESTMapper:             mgr.GetRESTMapper(),
		InsightsProxy:          insightsURL,
		ReconcilerTLS: common.NewReconcilerTLS(&common.ReconcilerTLSConfig{
			Client: mgr.GetClient(),
		}),
	}
}
