// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha2

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventInfo) DeepCopyInto(out *EventInfo) {
	*out = *in
	if in.Category != nil {
		in, out := &in.Category, &out.Category
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		*out = make(map[string]OptionDescriptor, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventInfo.
func (in *EventInfo) DeepCopy() *EventInfo {
	if in == nil {
		return nil
	}
	out := new(EventInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlightRecorder) DeepCopyInto(out *FlightRecorder) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlightRecorder.
func (in *FlightRecorder) DeepCopy() *FlightRecorder {
	if in == nil {
		return nil
	}
	out := new(FlightRecorder)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FlightRecorder) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlightRecorderList) DeepCopyInto(out *FlightRecorderList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FlightRecorder, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlightRecorderList.
func (in *FlightRecorderList) DeepCopy() *FlightRecorderList {
	if in == nil {
		return nil
	}
	out := new(FlightRecorderList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FlightRecorderList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlightRecorderSpec) DeepCopyInto(out *FlightRecorderSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlightRecorderSpec.
func (in *FlightRecorderSpec) DeepCopy() *FlightRecorderSpec {
	if in == nil {
		return nil
	}
	out := new(FlightRecorderSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlightRecorderStatus) DeepCopyInto(out *FlightRecorderStatus) {
	*out = *in
	if in.Events != nil {
		in, out := &in.Events, &out.Events
		*out = make([]EventInfo, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Target != nil {
		in, out := &in.Target, &out.Target
		*out = new(v1.ObjectReference)
		**out = **in
	}
	if in.RecordingSelector != nil {
		in, out := &in.RecordingSelector, &out.RecordingSelector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlightRecorderStatus.
func (in *FlightRecorderStatus) DeepCopy() *FlightRecorderStatus {
	if in == nil {
		return nil
	}
	out := new(FlightRecorderStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OptionDescriptor) DeepCopyInto(out *OptionDescriptor) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OptionDescriptor.
func (in *OptionDescriptor) DeepCopy() *OptionDescriptor {
	if in == nil {
		return nil
	}
	out := new(OptionDescriptor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Recording) DeepCopyInto(out *Recording) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Recording.
func (in *Recording) DeepCopy() *Recording {
	if in == nil {
		return nil
	}
	out := new(Recording)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Recording) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RecordingList) DeepCopyInto(out *RecordingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Recording, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RecordingList.
func (in *RecordingList) DeepCopy() *RecordingList {
	if in == nil {
		return nil
	}
	out := new(RecordingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RecordingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RecordingSpec) DeepCopyInto(out *RecordingSpec) {
	*out = *in
	if in.EventOptions != nil {
		in, out := &in.EventOptions, &out.EventOptions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.Duration = in.Duration
	if in.State != nil {
		in, out := &in.State, &out.State
		*out = new(RecordingState)
		**out = **in
	}
	if in.FlightRecorder != nil {
		in, out := &in.FlightRecorder, &out.FlightRecorder
		*out = new(v1.LocalObjectReference)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RecordingSpec.
func (in *RecordingSpec) DeepCopy() *RecordingSpec {
	if in == nil {
		return nil
	}
	out := new(RecordingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RecordingStatus) DeepCopyInto(out *RecordingStatus) {
	*out = *in
	if in.State != nil {
		in, out := &in.State, &out.State
		*out = new(RecordingState)
		**out = **in
	}
	in.StartTime.DeepCopyInto(&out.StartTime)
	out.Duration = in.Duration
	if in.DownloadURL != nil {
		in, out := &in.DownloadURL, &out.DownloadURL
		*out = new(string)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RecordingStatus.
func (in *RecordingStatus) DeepCopy() *RecordingStatus {
	if in == nil {
		return nil
	}
	out := new(RecordingStatus)
	in.DeepCopyInto(out)
	return out
}
