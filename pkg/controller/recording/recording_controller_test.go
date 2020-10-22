// Copyright (c) 2020 Red Hat, Inc.
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

package recording_test

import (
	"context"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	rhjmcv1beta1 "github.com/rh-jmc-team/container-jfr-operator/pkg/apis/rhjmc/v1beta1"
	jfrclient "github.com/rh-jmc-team/container-jfr-operator/pkg/client"
	"github.com/rh-jmc-team/container-jfr-operator/pkg/controller/recording"
	"github.com/rh-jmc-team/container-jfr-operator/test"
)

var _ = Describe("RecordingController", func() {
	var (
		objs       []runtime.Object
		handlers   []http.HandlerFunc
		server     *test.ContainerJFRServer
		client     client.Client
		controller *recording.ReconcileRecording
	)

	JustBeforeEach(func() {
		logf.SetLogger(zap.Logger())
		s := test.NewTestScheme()

		client = fake.NewFakeClientWithScheme(s, objs...)
		server = test.NewServer(client, handlers)
		controller = &recording.ReconcileRecording{
			Client:     client,
			Scheme:     s,
			Reconciler: test.NewTestReconciler(server, client),
		}
	})

	JustAfterEach(func() {
		server.Close()
	})

	BeforeEach(func() {
		objs = []runtime.Object{
			test.NewContainerJFR(), test.NewCACert(), test.NewFlightRecorder(),
			test.NewTargetPod(), test.NewContainerJFRService(), test.NewJMXAuthSecret(),
		}
	})

	AfterEach(func() {
		// Reset test inputs
		objs = nil
		handlers = nil
	})

	Describe("reconciling a request", func() {
		Context("with a new recording", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewRecording())
				handlers = []http.HandlerFunc{
					test.NewDumpHandler(),
					test.NewListHandler(test.NewRecordingDescriptors("RUNNING", 30000)),
				}
			})
			It("updates status with recording info", func() {
				desc := test.NewRecordingDescriptors("RUNNING", 30000)[0]
				expectRecordingStatus(controller, client, &desc)
			})
			It("adds finalizer to recording", func() {
				expectFinalizerPresent(controller, client)
			})
			It("should requeue after 10 seconds", func() {
				expectResult(controller, reconcile.Result{RequeueAfter: 10 * time.Second})
			})
		})
		Context("with a new recording that fails", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewRecording())
				handlers = []http.HandlerFunc{
					test.NewDumpFailHandler(),
				}
			})
			It("should requeue with error", func() {
				expectReconcileError(controller)
			})
		})
		Context("with a new continuous recording", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewContinuousRecording())
				handlers = []http.HandlerFunc{
					test.NewStartHandler(),
					test.NewListHandler(test.NewRecordingDescriptors("RUNNING", 0)),
				}
			})
			It("updates status with recording info", func() {
				desc := test.NewRecordingDescriptors("RUNNING", 0)[0]
				expectRecordingStatus(controller, client, &desc)
			})
			It("should requeue after 10 seconds", func() {
				expectResult(controller, reconcile.Result{RequeueAfter: 10 * time.Second})
			})
		})
		Context("with a new continuous recording that fails", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewContinuousRecording())
				handlers = []http.HandlerFunc{
					test.NewStartFailHandler(),
				}
			})
			It("should requeue with error", func() {
				expectReconcileError(controller)
			})
		})
		Context("with a running recording", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewRunningRecording())
				handlers = []http.HandlerFunc{
					test.NewListHandler(test.NewRecordingDescriptors("RUNNING", 30000)),
				}
			})
			It("should not change status", func() {
				expectStatusUnchanged(controller, client)
			})
			It("should requeue after 10 seconds", func() {
				expectResult(controller, reconcile.Result{RequeueAfter: 10 * time.Second})
			})
		})
		Context("with a running recording not found in Container JFR", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewRunningRecording())
				handlers = []http.HandlerFunc{
					test.NewListHandler([]jfrclient.RecordingDescriptor{}),
				}
			})
			It("should not change status", func() {
				expectStatusUnchanged(controller, client)
			})
			It("should requeue after 10 seconds", func() {
				expectResult(controller, reconcile.Result{RequeueAfter: 10 * time.Second})
			})
		})
		Context("when listing recordings fail", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewRunningRecording())
				handlers = []http.HandlerFunc{
					test.NewListFailHandler(test.NewRecordingDescriptors("RUNNING", 30000)),
				}
			})
			It("should requeue with error", func() {
				expectReconcileError(controller)
			})
		})
		Context("when listing recordings has unexpected state", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewRunningRecording())
				handlers = []http.HandlerFunc{
					test.NewListHandler(test.NewRecordingDescriptors("DOES-NOT-EXIST", 30000)),
				}
			})
			It("should requeue with error", func() {
				expectReconcileError(controller)
			})
		})
		Context("with a running recording to be stopped", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewRecordingToStop())
				handlers = []http.HandlerFunc{
					test.NewStopHandler(),
					test.NewListHandler(test.NewRecordingDescriptors("STOPPED", 0)),
				}
			})
			It("should stop recording", func() {
				desc := test.NewRecordingDescriptors("STOPPED", 0)[0]
				expectRecordingStatus(controller, client, &desc)
			})
			It("should not requeue", func() {
				expectResult(controller, reconcile.Result{})
			})
		})
		Context("with a running recording to be stopped that fails", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewRecordingToStop())
				handlers = []http.HandlerFunc{
					test.NewStopFailHandler(),
				}
			})
			It("should requeue with error", func() {
				expectReconcileError(controller)
			})
		})
		Context("with a stopped recording to be archived", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewStoppedRecordingToArchive())
				handlers = []http.HandlerFunc{
					test.NewListHandler(test.NewRecordingDescriptors("STOPPED", 30000)),
					test.NewListSavedHandler([]jfrclient.SavedRecording{}),
					test.NewSaveHandler(),
					test.NewListSavedHandler(test.NewSavedRecordings()),
				}
			})
			It("should update download URL", func() {
				req := reconcile.Request{NamespacedName: types.NamespacedName{Name: "my-recording", Namespace: "default"}}

				before := &rhjmcv1beta1.Recording{}
				err := client.Get(context.Background(), req.NamespacedName, before)
				Expect(err).ToNot(HaveOccurred())

				_, err = controller.Reconcile(req)
				Expect(err).ToNot(HaveOccurred())

				after := &rhjmcv1beta1.Recording{}
				err = client.Get(context.Background(), req.NamespacedName, after)
				Expect(err).ToNot(HaveOccurred())

				// Should all be the same except for Download URL
				saved := test.NewSavedRecordings()[0]
				Expect(after.Status.State).To(Equal(before.Status.State))
				Expect(after.Status.Duration).To(Equal(before.Status.Duration))
				Expect(after.Status.StartTime).To(Equal(before.Status.StartTime))
				Expect(after.Status.DownloadURL).ToNot(BeNil())
				Expect(*after.Status.DownloadURL).To(Equal(saved.DownloadURL))
			})
			It("should not requeue", func() {
				expectResult(controller, reconcile.Result{})
			})
		})
		Context("when listing saved recordings fails", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewStoppedRecordingToArchive())
				handlers = []http.HandlerFunc{
					test.NewListHandler(test.NewRecordingDescriptors("STOPPED", 30000)),
					test.NewListSavedFailHandler(test.NewSavedRecordings()),
				}
			})
			It("should requeue with error", func() {
				expectReconcileError(controller)
			})
		})
		Context("when archiving recording fails", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewStoppedRecordingToArchive())
				handlers = []http.HandlerFunc{
					test.NewListHandler(test.NewRecordingDescriptors("STOPPED", 30000)),
					test.NewListSavedHandler(test.NewSavedRecordings()),
					test.NewSaveFailHandler(),
				}
			})
			It("should requeue with error", func() {
				expectReconcileError(controller)
			})
		})
		Context("with a running recording to be stopped and archived", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewRecordingToStopAndArchive())
				handlers = []http.HandlerFunc{
					test.NewStopHandler(),
					test.NewListHandler(test.NewRecordingDescriptors("STOPPED", 30000)),
					test.NewListSavedHandler([]jfrclient.SavedRecording{}),
					test.NewSaveHandler(),
					test.NewListSavedHandler(test.NewSavedRecordings()),
				}
			})
			It("should stop recording", func() {
				req := reconcile.Request{NamespacedName: types.NamespacedName{Name: "my-recording", Namespace: "default"}}
				_, err := controller.Reconcile(req)
				Expect(err).ToNot(HaveOccurred())

				obj := &rhjmcv1beta1.Recording{}
				err = client.Get(context.Background(), req.NamespacedName, obj)
				Expect(err).ToNot(HaveOccurred())

				Expect(obj.Status.State).ToNot(BeNil())
				Expect(*obj.Status.State).To(Equal(rhjmcv1beta1.RecordingStateStopped))
			})
			It("should update download URL", func() {
				req := reconcile.Request{NamespacedName: types.NamespacedName{Name: "my-recording", Namespace: "default"}}
				_, err := controller.Reconcile(req)
				Expect(err).ToNot(HaveOccurred())

				obj := &rhjmcv1beta1.Recording{}
				err = client.Get(context.Background(), req.NamespacedName, obj)
				Expect(err).ToNot(HaveOccurred())

				Expect(obj.Status.DownloadURL).ToNot(BeNil())
				Expect(*obj.Status.DownloadURL).To(Equal("http://path/to/saved-test-recording.jfr"))
			})
			It("should not requeue", func() {
				expectResult(controller, reconcile.Result{})
			})
		})
		Context("with an archived recording", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewArchivedRecording())
				handlers = []http.HandlerFunc{
					test.NewListHandler(test.NewRecordingDescriptors("STOPPED", 30000)),
					test.NewListSavedHandler(test.NewSavedRecordings()),
				}
			})
			It("should not change status", func() {
				expectStatusUnchanged(controller, client)
			})
			It("should not requeue", func() {
				expectResult(controller, reconcile.Result{})
			})
		})
		Context("with a deleted archived recording", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewDeletedArchivedRecording())
				handlers = []http.HandlerFunc{
					test.NewListSavedHandler(test.NewSavedRecordings()),
					test.NewDeleteSavedHandler(),
					test.NewListHandler(test.NewRecordingDescriptors("STOPPED", 30000)),
					test.NewDeleteHandler(),
				}
			})
			It("should remove the finalizer", func() {
				expectFinalizerAbsent(controller, client)
			})
			It("should not requeue", func() {
				expectResult(controller, reconcile.Result{})
			})
		})
		Context("with a deleted recording with missing FlightRecorder", func() {
			BeforeEach(func() {
				objs = []runtime.Object{
					test.NewContainerJFR(), test.NewCACert(), test.NewTargetPod(),
					test.NewContainerJFRService(), test.NewDeletedArchivedRecording(),
					test.NewJMXAuthSecret(),
				}
				handlers = []http.HandlerFunc{
					test.NewListSavedHandler(test.NewSavedRecordings()),
					test.NewDeleteSavedHandler(),
				}
			})
			It("should remove the finalizer", func() {
				expectFinalizerAbsent(controller, client)
			})
			It("should not requeue", func() {
				expectResult(controller, reconcile.Result{})
			})
		})
		Context("when deleting the saved recording fails", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewDeletedArchivedRecording())
				handlers = []http.HandlerFunc{
					test.NewListSavedHandler(test.NewSavedRecordings()),
					test.NewDeleteSavedFailHandler(),
				}
			})
			It("should keep the finalizer", func() {
				expectFinalizerPresent(controller, client)
			})
			It("should requeue with error", func() {
				expectReconcileError(controller)
			})
		})
		Context("when deleting the in-memory recording fails", func() {
			BeforeEach(func() {
				objs = append(objs, test.NewDeletedArchivedRecording())
				handlers = []http.HandlerFunc{
					test.NewListSavedHandler(test.NewSavedRecordings()),
					test.NewDeleteSavedHandler(),
					test.NewListHandler(test.NewRecordingDescriptors("STOPPED", 30000)),
					test.NewDeleteFailHandler(),
				}
			})
			It("should keep the finalizer", func() {
				expectFinalizerPresent(controller, client)
			})
			It("should requeue with error", func() {
				expectReconcileError(controller)
			})
		})
		Context("Recording does not exist", func() {
			BeforeEach(func() {
				handlers = []http.HandlerFunc{}
			})
			It("should do nothing", func() {
				req := reconcile.Request{NamespacedName: types.NamespacedName{Name: "does-not-exist", Namespace: "default"}}
				result, err := controller.Reconcile(req)
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(reconcile.Result{}))
			})
		})
		Context("FlightRecorder Status not updated yet", func() {
			BeforeEach(func() {
				otherFr := test.NewFlightRecorder()
				otherFr.Status = rhjmcv1beta1.FlightRecorderStatus{}
				objs = []runtime.Object{
					test.NewContainerJFR(), test.NewCACert(), otherFr, test.NewTargetPod(),
					test.NewContainerJFRService(), test.NewRecording(), test.NewJMXAuthSecret(),
				}
				handlers = []http.HandlerFunc{}
			})
			It("should requeue", func() {
				expectResult(controller, reconcile.Result{RequeueAfter: time.Second})
			})
		})
		Context("Container JFR CR is missing", func() {
			BeforeEach(func() {
				objs = []runtime.Object{
					test.NewFlightRecorder(), test.NewCACert(), test.NewTargetPod(),
					test.NewContainerJFRService(), test.NewRecording(), test.NewJMXAuthSecret(),
				}
				handlers = []http.HandlerFunc{}
			})
			It("should requeue with error", func() {
				expectReconcileError(controller)
			})
		})
		Context("Container JFR service is missing", func() {
			BeforeEach(func() {
				objs = []runtime.Object{
					test.NewContainerJFR(), test.NewCACert(), test.NewFlightRecorder(),
					test.NewTargetPod(), test.NewRecording(), test.NewJMXAuthSecret(),
				}
				handlers = []http.HandlerFunc{}
			})
			It("should requeue with error", func() {
				expectReconcileError(controller)
			})
		})
		Context("FlightRecorder is missing", func() {
			BeforeEach(func() {
				objs = []runtime.Object{
					test.NewContainerJFR(), test.NewCACert(), test.NewTargetPod(),
					test.NewContainerJFRService(), test.NewRecording(), test.NewJMXAuthSecret(),
				}
				handlers = []http.HandlerFunc{}
			})
			It("should not requeue", func() {
				req := reconcile.Request{NamespacedName: types.NamespacedName{Name: "my-recording", Namespace: "default"}}
				result, err := controller.Reconcile(req)
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(reconcile.Result{}))
			})
		})
		Context("Target pod is missing", func() {
			BeforeEach(func() {
				objs = []runtime.Object{
					test.NewContainerJFR(), test.NewCACert(), test.NewFlightRecorder(),
					test.NewContainerJFRService(), test.NewRecording(), test.NewJMXAuthSecret(),
				}
				handlers = []http.HandlerFunc{}
			})
			It("should requeue with error", func() {
				expectReconcileError(controller)
			})
		})
		Context("Target pod has no IP", func() {
			BeforeEach(func() {
				otherPod := test.NewTargetPod()
				otherPod.Status.PodIP = ""
				objs = []runtime.Object{
					test.NewContainerJFR(), test.NewCACert(), test.NewFlightRecorder(), otherPod,
					test.NewContainerJFRService(), test.NewRecording(), test.NewJMXAuthSecret(),
				}
				handlers = []http.HandlerFunc{}
			})
			It("should requeue with error", func() {
				expectReconcileError(controller)
			})
		})
	})
})

func expectRecordingStatus(controller *recording.ReconcileRecording, client client.Client, desc *jfrclient.RecordingDescriptor) {
	obj := reconcileAndGet(controller, client)

	Expect(obj.Status.State).ToNot(BeNil())
	Expect(*obj.Status.State).To(Equal(rhjmcv1beta1.RecordingState(desc.State)))
	// Converted to RFC3339 during serialization (sub-second precision lost)
	Expect(obj.Status.StartTime).To(Equal(metav1.Unix(0, desc.StartTime*int64(time.Millisecond)).Rfc3339Copy()))
	Expect(obj.Status.Duration).To(Equal(metav1.Duration{
		Duration: time.Duration(desc.Duration) * time.Millisecond,
	}))
	Expect(obj.Status.DownloadURL).ToNot(BeNil())
	Expect(*obj.Status.DownloadURL).To(Equal(desc.DownloadURL))
}

func expectStatusUnchanged(controller *recording.ReconcileRecording, client client.Client) {
	before := &rhjmcv1beta1.Recording{}
	err := client.Get(context.Background(), types.NamespacedName{Name: "my-recording", Namespace: "default"}, before)
	Expect(err).ToNot(HaveOccurred())

	after := reconcileAndGet(controller, client)
	Expect(after.Status).To(Equal(before.Status))
}

func expectFinalizerPresent(controller *recording.ReconcileRecording, client client.Client) {
	obj := reconcileAndGet(controller, client)
	finalizers := obj.GetFinalizers()
	Expect(finalizers).To(ContainElement("recording.finalizer.rhjmc.redhat.com"))
}

func expectFinalizerAbsent(controller *recording.ReconcileRecording, client client.Client) {
	obj := reconcileAndGet(controller, client)
	finalizers := obj.GetFinalizers()
	Expect(finalizers).ToNot(ContainElement("recording.finalizer.rhjmc.redhat.com"))
}

func expectReconcileError(controller *recording.ReconcileRecording) {
	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: "my-recording", Namespace: "default"}}
	result, err := controller.Reconcile(req)
	Expect(err).To(HaveOccurred())
	Expect(result).To(Equal(reconcile.Result{}))
}

func expectResult(controller *recording.ReconcileRecording, result reconcile.Result) {
	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: "my-recording", Namespace: "default"}}
	result, err := controller.Reconcile(req)
	Expect(err).ToNot(HaveOccurred())
	Expect(result).To(Equal(result))
}

func reconcileAndGet(controller *recording.ReconcileRecording, client client.Client) *rhjmcv1beta1.Recording {
	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: "my-recording", Namespace: "default"}}
	controller.Reconcile(req)

	obj := &rhjmcv1beta1.Recording{}
	err := client.Get(context.Background(), req.NamespacedName, obj)
	Expect(err).ToNot(HaveOccurred())
	return obj
}
