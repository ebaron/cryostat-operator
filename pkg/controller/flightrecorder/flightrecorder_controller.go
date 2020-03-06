package flightrecorder

import (
	"context"
	"time"

	rhjmcv1alpha2 "github.com/rh-jmc-team/container-jfr-operator/pkg/apis/rhjmc/v1alpha2"
	jfrclient "github.com/rh-jmc-team/container-jfr-operator/pkg/client"
	common "github.com/rh-jmc-team/container-jfr-operator/pkg/controller/common"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_flightrecorder")

// Add creates a new FlightRecorder Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileFlightRecorder{scheme: mgr.GetScheme(),
		CommonReconciler: &common.CommonReconciler{
			Client: mgr.GetClient(),
		},
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("flightrecorder-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource FlightRecorder
	err = c.Watch(&source.Kind{Type: &rhjmcv1alpha2.FlightRecorder{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileFlightRecorder implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileFlightRecorder{}

// ReconcileFlightRecorder reconciles a FlightRecorder object
type ReconcileFlightRecorder struct {
	scheme *runtime.Scheme
	*common.CommonReconciler
}

// Reconcile reads that state of the cluster for a FlightRecorder object and makes changes based on the state read
// and what is in the FlightRecorder.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileFlightRecorder) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling FlightRecorder")

	cjfr, err := r.FindContainerJFR(ctx, request.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Keep client open to Container JFR as long as it doesn't fail
	if r.JfrClient == nil {
		jfrClient, err := r.ConnectToContainerJFR(ctx, cjfr.Namespace, cjfr.Name)
		if err != nil {
			// Need Container JFR in order to reconcile anything, requeue until it appears
			return reconcile.Result{}, err
		}
		r.JfrClient = jfrClient
	}

	// Fetch the FlightRecorder instance
	instance := &rhjmcv1alpha2.FlightRecorder{}
	err = r.Client.Get(ctx, request.NamespacedName, instance)
	if err != nil {
		if kerrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Look up service corresponding to this FlightRecorder object
	targetRef := instance.Status.Target
	if targetRef == nil {
		// FlightRecorder status must not have been updated yet
		return reconcile.Result{RequeueAfter: time.Second}, nil
	}
	targetSvc := &corev1.Service{}
	err = r.Client.Get(ctx, types.NamespacedName{Namespace: targetRef.Namespace, Name: targetRef.Name}, targetSvc)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Tell Container JFR to connect to the target service
	jfrclient.ClientLock.Lock()
	defer jfrclient.ClientLock.Unlock()
	err = r.ConnectToService(targetSvc, instance.Status.Port)
	if err != nil {
		return reconcile.Result{}, err
	}
	defer r.DisconnectClient()

	// FIXME Having two clients connecting to potentially different services could cause problems. Maybe need to synchronize with one client
	/*
		// Check for any new recording requests in this FlightRecorder's spec
			// and instruct Container JFR to create corresponding recordings
			log.Info("Syncing recording requests for service", "service", targetSvc.Name, "namespace", targetSvc.Namespace,
				"host", *clusterIP, "port", jmxPort)
		f	or _, request := range instance.Spec.RecordingRequests {
				log.Info("Creating new recording", "name", request.Name, "duration", request.Duration, "eventOptions", request.EventOptions)
				err := r.jfrClient.DumpRecording(request.Name, int(request.Duration.Seconds()), request.EventOptions)
				if err != nil {
					log.Error(err, "failed to create new recording")
					r.closeClient() // TODO maybe track an error state in the client instead of relying on calling this everywhere
					return reconcile.Result{}, err
				}
		}

		// Get an updated list of in-memory flight recordings
			log.Info("Listing recordings for service", "service", targetSvc.Name, "namespace", targetSvc.Namespace,
				"host", *clusterIP, "port", jmxPort)
		d	escriptors, err := r.jfrClient.ListRecordings()
			if err != nil {
				log.Error(err, "failed to list flight recordings")
				r.closeClient()
				return reconcile.Result{}, err
		}

		reqLogger.Info("Updating FlightRecorder", "Namespace", instance.Namespace, "Name", instance.Name)
			// Remove any recording requests from the spec that are now showing in Container JFR's list
			newRequests := []rhjmcv1alpha1.RecordingRequest{}
			for _, req := range instance.Spec.RecordingRequests {
				for _, desc := range descriptors {
					if req.Name != desc.Name {
						newRequests = append(newRequests, req)
						break
				}
				}
		}
			instance.Spec.RecordingRequests = newRequests
			err = r.client.Update(ctx, instance)
			if err != nil {
				return reconcile.Result{}, err
		}

		// Update recording info in Status with info received from Container JFR
			recordings := createRecordingInfo(descriptors)

		// TODO Download URLs returned by Container JFR's 'list' command currently
			// work when it is connected to the target JVM. To work around this,
			// we archive the recording to persistent storage and update the download
			// URL to point to that saved file.
			toUpdate := map[string]*rhjmcv1alpha1.RecordingInfo{}
			for idx, newInfo := range recordings {
				oldInfo := findRecordingByName(instance.Status.Recordings, newInfo.Name)
				// Recording completed since last observation
				if !newInfo.Active && (oldInfo == nil || oldInfo.Active) {
					filename, err := r.jfrClient.SaveRecording(newInfo.Name)
					if err != nil {
						log.Error(err, "failed to save recording", "name", newInfo.Name)
						return reconcile.Result{}, err
				}
					toUpdate[*filename] = &recordings[idx]
				} else if oldInfo != nil && len(oldInfo.DownloadURL) > 0 {
					// Use previously obtained download URL
					recordings[idx].DownloadURL = oldInfo.DownloadURL
				}
		}

		if len(toUpdate) > 0 {
				savedRecordings, err := r.jfrClient.ListSavedRecordings()
				if err != nil {
					return reconcile.Result{}, err
				}
				// Update download URLs using list of saved recordings
				for _, saved := range savedRecordings {
					if info, pres := toUpdate[saved.Name]; pres {
						log.Info("updating download URL", "name", info.Name, "url", saved.DownloadURL)
						info.DownloadURL = saved.DownloadURL
				}
			}
			}

		instance.Status.Recordings = recordings
			err = r.client.Status().Update(ctx, instance)
			if err != nil {
				return reconcile.Result{}, err
			}

		// Requeue if any recordings are still in progress
			result := reconcile.Result{}
			for _, recording := range recordings {
				if recording.Active {
					// Check progress of recordings after 10 seconds
					result = reconcile.Result{RequeueAfter: 10 * time.Second}
					break
			}
		}
	*/

	reqLogger.Info("FlightRecorder successfully updated", "Namespace", instance.Namespace, "Name", instance.Name)
	return reconcile.Result{}, nil
}

/*func findRecordingByName(recordings []rhjmcv1alpha1.RecordingInfo, name string) *rhjmcv1alpha1.RecordingInfo {
	for idx, recording := range recordings {
		if recording.Name == name {
			return &recordings[idx]
		}
	}
	return nil
}

func createRecordingInfo(descriptors []jfrclient.RecordingDescriptor) []rhjmcv1alpha1.RecordingInfo {
	infos := make([]rhjmcv1alpha1.RecordingInfo, len(descriptors))
	for i, descriptor := range descriptors {
		// Consider any recording not stopped to be "active"
		active := descriptor.State != jfrclient.RecordingStateStopped
		startTime := metav1.Unix(0, descriptor.StartTime*int64(time.Millisecond))
		duration := metav1.Duration{
			Duration: time.Duration(descriptor.Duration) * time.Millisecond,
		}
		info := rhjmcv1alpha1.RecordingInfo{
			Name:      descriptor.Name,
			Active:    active,
			StartTime: startTime,
			Duration:  duration,
		}
		infos[i] = info
	}
	return infos
}*/
