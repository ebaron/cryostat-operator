//go:build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Affinity) DeepCopyInto(out *Affinity) {
	*out = *in
	if in.NodeAffinity != nil {
		in, out := &in.NodeAffinity, &out.NodeAffinity
		*out = new(corev1.NodeAffinity)
		(*in).DeepCopyInto(*out)
	}
	if in.PodAffinity != nil {
		in, out := &in.PodAffinity, &out.PodAffinity
		*out = new(corev1.PodAffinity)
		(*in).DeepCopyInto(*out)
	}
	if in.PodAntiAffinity != nil {
		in, out := &in.PodAntiAffinity, &out.PodAntiAffinity
		*out = new(corev1.PodAntiAffinity)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Affinity.
func (in *Affinity) DeepCopy() *Affinity {
	if in == nil {
		return nil
	}
	out := new(Affinity)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AuthorizationProperties) DeepCopyInto(out *AuthorizationProperties) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthorizationProperties.
func (in *AuthorizationProperties) DeepCopy() *AuthorizationProperties {
	if in == nil {
		return nil
	}
	out := new(AuthorizationProperties)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CertificateSecret) DeepCopyInto(out *CertificateSecret) {
	*out = *in
	if in.CertificateKey != nil {
		in, out := &in.CertificateKey, &out.CertificateKey
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CertificateSecret.
func (in *CertificateSecret) DeepCopy() *CertificateSecret {
	if in == nil {
		return nil
	}
	out := new(CertificateSecret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CoreServiceConfig) DeepCopyInto(out *CoreServiceConfig) {
	*out = *in
	if in.HTTPPort != nil {
		in, out := &in.HTTPPort, &out.HTTPPort
		*out = new(int32)
		**out = **in
	}
	if in.JMXPort != nil {
		in, out := &in.JMXPort, &out.JMXPort
		*out = new(int32)
		**out = **in
	}
	in.ServiceConfig.DeepCopyInto(&out.ServiceConfig)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CoreServiceConfig.
func (in *CoreServiceConfig) DeepCopy() *CoreServiceConfig {
	if in == nil {
		return nil
	}
	out := new(CoreServiceConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Cryostat) DeepCopyInto(out *Cryostat) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Cryostat.
func (in *Cryostat) DeepCopy() *Cryostat {
	if in == nil {
		return nil
	}
	out := new(Cryostat)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Cryostat) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CryostatList) DeepCopyInto(out *CryostatList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Cryostat, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CryostatList.
func (in *CryostatList) DeepCopy() *CryostatList {
	if in == nil {
		return nil
	}
	out := new(CryostatList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CryostatList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CryostatSpec) DeepCopyInto(out *CryostatSpec) {
	*out = *in
	if in.TrustedCertSecrets != nil {
		in, out := &in.TrustedCertSecrets, &out.TrustedCertSecrets
		*out = make([]CertificateSecret, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.EventTemplates != nil {
		in, out := &in.EventTemplates, &out.EventTemplates
		*out = make([]TemplateConfigMap, len(*in))
		copy(*out, *in)
	}
	if in.EnableCertManager != nil {
		in, out := &in.EnableCertManager, &out.EnableCertManager
		*out = new(bool)
		**out = **in
	}
	if in.StorageOptions != nil {
		in, out := &in.StorageOptions, &out.StorageOptions
		*out = new(StorageConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.ServiceOptions != nil {
		in, out := &in.ServiceOptions, &out.ServiceOptions
		*out = new(ServiceConfigList)
		(*in).DeepCopyInto(*out)
	}
	if in.NetworkOptions != nil {
		in, out := &in.NetworkOptions, &out.NetworkOptions
		*out = new(NetworkConfigurationList)
		(*in).DeepCopyInto(*out)
	}
	if in.ReportOptions != nil {
		in, out := &in.ReportOptions, &out.ReportOptions
		*out = new(ReportConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.JmxCacheOptions != nil {
		in, out := &in.JmxCacheOptions, &out.JmxCacheOptions
		*out = new(JmxCacheOptions)
		**out = **in
	}
	in.Resources.DeepCopyInto(&out.Resources)
	if in.AuthProperties != nil {
		in, out := &in.AuthProperties, &out.AuthProperties
		*out = new(AuthorizationProperties)
		**out = **in
	}
	if in.SecurityOptions != nil {
		in, out := &in.SecurityOptions, &out.SecurityOptions
		*out = new(SecurityOptions)
		(*in).DeepCopyInto(*out)
	}
	if in.SchedulingOptions != nil {
		in, out := &in.SchedulingOptions, &out.SchedulingOptions
		*out = new(SchedulingConfiguration)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CryostatSpec.
func (in *CryostatSpec) DeepCopy() *CryostatSpec {
	if in == nil {
		return nil
	}
	out := new(CryostatSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CryostatStatus) DeepCopyInto(out *CryostatStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CryostatStatus.
func (in *CryostatStatus) DeepCopy() *CryostatStatus {
	if in == nil {
		return nil
	}
	out := new(CryostatStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EmptyDirConfig) DeepCopyInto(out *EmptyDirConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EmptyDirConfig.
func (in *EmptyDirConfig) DeepCopy() *EmptyDirConfig {
	if in == nil {
		return nil
	}
	out := new(EmptyDirConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GrafanaServiceConfig) DeepCopyInto(out *GrafanaServiceConfig) {
	*out = *in
	if in.HTTPPort != nil {
		in, out := &in.HTTPPort, &out.HTTPPort
		*out = new(int32)
		**out = **in
	}
	in.ServiceConfig.DeepCopyInto(&out.ServiceConfig)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GrafanaServiceConfig.
func (in *GrafanaServiceConfig) DeepCopy() *GrafanaServiceConfig {
	if in == nil {
		return nil
	}
	out := new(GrafanaServiceConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JmxCacheOptions) DeepCopyInto(out *JmxCacheOptions) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JmxCacheOptions.
func (in *JmxCacheOptions) DeepCopy() *JmxCacheOptions {
	if in == nil {
		return nil
	}
	out := new(JmxCacheOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkConfiguration) DeepCopyInto(out *NetworkConfiguration) {
	*out = *in
	if in.IngressSpec != nil {
		in, out := &in.IngressSpec, &out.IngressSpec
		*out = new(networkingv1.IngressSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkConfiguration.
func (in *NetworkConfiguration) DeepCopy() *NetworkConfiguration {
	if in == nil {
		return nil
	}
	out := new(NetworkConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkConfigurationList) DeepCopyInto(out *NetworkConfigurationList) {
	*out = *in
	if in.CoreConfig != nil {
		in, out := &in.CoreConfig, &out.CoreConfig
		*out = new(NetworkConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.CommandConfig != nil {
		in, out := &in.CommandConfig, &out.CommandConfig
		*out = new(NetworkConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.GrafanaConfig != nil {
		in, out := &in.GrafanaConfig, &out.GrafanaConfig
		*out = new(NetworkConfiguration)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkConfigurationList.
func (in *NetworkConfigurationList) DeepCopy() *NetworkConfigurationList {
	if in == nil {
		return nil
	}
	out := new(NetworkConfigurationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PersistentVolumeClaimConfig) DeepCopyInto(out *PersistentVolumeClaimConfig) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(corev1.PersistentVolumeClaimSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PersistentVolumeClaimConfig.
func (in *PersistentVolumeClaimConfig) DeepCopy() *PersistentVolumeClaimConfig {
	if in == nil {
		return nil
	}
	out := new(PersistentVolumeClaimConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReportConfiguration) DeepCopyInto(out *ReportConfiguration) {
	*out = *in
	in.Resources.DeepCopyInto(&out.Resources)
	if in.SecurityOptions != nil {
		in, out := &in.SecurityOptions, &out.SecurityOptions
		*out = new(ReportsSecurityOptions)
		(*in).DeepCopyInto(*out)
	}
	if in.SchedulingOptions != nil {
		in, out := &in.SchedulingOptions, &out.SchedulingOptions
		*out = new(SchedulingConfiguration)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReportConfiguration.
func (in *ReportConfiguration) DeepCopy() *ReportConfiguration {
	if in == nil {
		return nil
	}
	out := new(ReportConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReportsSecurityOptions) DeepCopyInto(out *ReportsSecurityOptions) {
	*out = *in
	if in.PodSecurityContext != nil {
		in, out := &in.PodSecurityContext, &out.PodSecurityContext
		*out = new(corev1.PodSecurityContext)
		(*in).DeepCopyInto(*out)
	}
	if in.ReportsSecurityContext != nil {
		in, out := &in.ReportsSecurityContext, &out.ReportsSecurityContext
		*out = new(corev1.SecurityContext)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReportsSecurityOptions.
func (in *ReportsSecurityOptions) DeepCopy() *ReportsSecurityOptions {
	if in == nil {
		return nil
	}
	out := new(ReportsSecurityOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReportsServiceConfig) DeepCopyInto(out *ReportsServiceConfig) {
	*out = *in
	if in.HTTPPort != nil {
		in, out := &in.HTTPPort, &out.HTTPPort
		*out = new(int32)
		**out = **in
	}
	in.ServiceConfig.DeepCopyInto(&out.ServiceConfig)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReportsServiceConfig.
func (in *ReportsServiceConfig) DeepCopy() *ReportsServiceConfig {
	if in == nil {
		return nil
	}
	out := new(ReportsServiceConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceConfigList) DeepCopyInto(out *ResourceConfigList) {
	*out = *in
	in.CoreResources.DeepCopyInto(&out.CoreResources)
	in.DataSourceResources.DeepCopyInto(&out.DataSourceResources)
	in.GrafanaResources.DeepCopyInto(&out.GrafanaResources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceConfigList.
func (in *ResourceConfigList) DeepCopy() *ResourceConfigList {
	if in == nil {
		return nil
	}
	out := new(ResourceConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SchedulingConfiguration) DeepCopyInto(out *SchedulingConfiguration) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]corev1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SchedulingConfiguration.
func (in *SchedulingConfiguration) DeepCopy() *SchedulingConfiguration {
	if in == nil {
		return nil
	}
	out := new(SchedulingConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecurityOptions) DeepCopyInto(out *SecurityOptions) {
	*out = *in
	if in.PodSecurityContext != nil {
		in, out := &in.PodSecurityContext, &out.PodSecurityContext
		*out = new(corev1.PodSecurityContext)
		(*in).DeepCopyInto(*out)
	}
	if in.CoreSecurityContext != nil {
		in, out := &in.CoreSecurityContext, &out.CoreSecurityContext
		*out = new(corev1.SecurityContext)
		(*in).DeepCopyInto(*out)
	}
	if in.DataSourceSecurityContext != nil {
		in, out := &in.DataSourceSecurityContext, &out.DataSourceSecurityContext
		*out = new(corev1.SecurityContext)
		(*in).DeepCopyInto(*out)
	}
	if in.GrafanaSecurityContext != nil {
		in, out := &in.GrafanaSecurityContext, &out.GrafanaSecurityContext
		*out = new(corev1.SecurityContext)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecurityOptions.
func (in *SecurityOptions) DeepCopy() *SecurityOptions {
	if in == nil {
		return nil
	}
	out := new(SecurityOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceConfig) DeepCopyInto(out *ServiceConfig) {
	*out = *in
	if in.ServiceType != nil {
		in, out := &in.ServiceType, &out.ServiceType
		*out = new(corev1.ServiceType)
		**out = **in
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceConfig.
func (in *ServiceConfig) DeepCopy() *ServiceConfig {
	if in == nil {
		return nil
	}
	out := new(ServiceConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceConfigList) DeepCopyInto(out *ServiceConfigList) {
	*out = *in
	if in.CoreConfig != nil {
		in, out := &in.CoreConfig, &out.CoreConfig
		*out = new(CoreServiceConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.GrafanaConfig != nil {
		in, out := &in.GrafanaConfig, &out.GrafanaConfig
		*out = new(GrafanaServiceConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.ReportsConfig != nil {
		in, out := &in.ReportsConfig, &out.ReportsConfig
		*out = new(ReportsServiceConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceConfigList.
func (in *ServiceConfigList) DeepCopy() *ServiceConfigList {
	if in == nil {
		return nil
	}
	out := new(ServiceConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageConfiguration) DeepCopyInto(out *StorageConfiguration) {
	*out = *in
	if in.PVC != nil {
		in, out := &in.PVC, &out.PVC
		*out = new(PersistentVolumeClaimConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.EmptyDir != nil {
		in, out := &in.EmptyDir, &out.EmptyDir
		*out = new(EmptyDirConfig)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageConfiguration.
func (in *StorageConfiguration) DeepCopy() *StorageConfiguration {
	if in == nil {
		return nil
	}
	out := new(StorageConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TemplateConfigMap) DeepCopyInto(out *TemplateConfigMap) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TemplateConfigMap.
func (in *TemplateConfigMap) DeepCopy() *TemplateConfigMap {
	if in == nil {
		return nil
	}
	out := new(TemplateConfigMap)
	in.DeepCopyInto(out)
	return out
}
