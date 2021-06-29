// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Autoscaling) DeepCopyInto(out *Autoscaling) {
	*out = *in
	if in.MinReplicas != nil {
		in, out := &in.MinReplicas, &out.MinReplicas
		*out = new(int32)
		**out = **in
	}
	if in.AverageCPUUtilization != nil {
		in, out := &in.AverageCPUUtilization, &out.AverageCPUUtilization
		*out = new(int32)
		**out = **in
	}
	if in.ScaleDownStabilizationWindowSeconds != nil {
		in, out := &in.ScaleDownStabilizationWindowSeconds, &out.ScaleDownStabilizationWindowSeconds
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Autoscaling.
func (in *Autoscaling) DeepCopy() *Autoscaling {
	if in == nil {
		return nil
	}
	out := new(Autoscaling)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterConfig) DeepCopyInto(out *ClusterConfig) {
	*out = *in
	out.IstioConfig = in.IstioConfig
	if in.GlobalLabels != nil {
		in, out := &in.GlobalLabels, &out.GlobalLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Image != nil {
		in, out := &in.Image, &out.Image
		*out = new(OCIImageDefinition)
		**out = **in
	}
	if in.Autoscaling != nil {
		in, out := &in.Autoscaling, &out.Autoscaling
		*out = new(Autoscaling)
		(*in).DeepCopyInto(*out)
	}
	in.NetworkPolicy.DeepCopyInto(&out.NetworkPolicy)
	out.ServiceAccount = in.ServiceAccount
	if in.KerberosKeytab != nil {
		in, out := &in.KerberosKeytab, &out.KerberosKeytab
		*out = new(KerberosKeytabConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.ImagePullSecrets != nil {
		in, out := &in.ImagePullSecrets, &out.ImagePullSecrets
		*out = make([]v1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
	if in.PodSecurityContext != nil {
		in, out := &in.PodSecurityContext, &out.PodSecurityContext
		*out = new(v1.PodSecurityContext)
		(*in).DeepCopyInto(*out)
	}
	if in.EnvVars != nil {
		in, out := &in.EnvVars, &out.EnvVars
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterConfig.
func (in *ClusterConfig) DeepCopy() *ClusterConfig {
	if in == nil {
		return nil
	}
	out := new(ClusterConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterStatusConfig) DeepCopyInto(out *ClusterStatusConfig) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterStatusConfig.
func (in *ClusterStatusConfig) DeepCopy() *ClusterStatusConfig {
	if in == nil {
		return nil
	}
	out := new(ClusterStatusConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DaskCluster) DeepCopyInto(out *DaskCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DaskCluster.
func (in *DaskCluster) DeepCopy() *DaskCluster {
	if in == nil {
		return nil
	}
	out := new(DaskCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DaskCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DaskClusterList) DeepCopyInto(out *DaskClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DaskCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DaskClusterList.
func (in *DaskClusterList) DeepCopy() *DaskClusterList {
	if in == nil {
		return nil
	}
	out := new(DaskClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DaskClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DaskClusterSpec) DeepCopyInto(out *DaskClusterSpec) {
	*out = *in
	in.ClusterConfig.DeepCopyInto(&out.ClusterConfig)
	in.Scheduler.DeepCopyInto(&out.Scheduler)
	in.Worker.DeepCopyInto(&out.Worker)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DaskClusterSpec.
func (in *DaskClusterSpec) DeepCopy() *DaskClusterSpec {
	if in == nil {
		return nil
	}
	out := new(DaskClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DaskClusterStatus) DeepCopyInto(out *DaskClusterStatus) {
	*out = *in
	in.ClusterStatusConfig.DeepCopyInto(&out.ClusterStatusConfig)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DaskClusterStatus.
func (in *DaskClusterStatus) DeepCopy() *DaskClusterStatus {
	if in == nil {
		return nil
	}
	out := new(DaskClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DaskClusterWorker) DeepCopyInto(out *DaskClusterWorker) {
	*out = *in
	in.WorkloadConfig.DeepCopyInto(&out.WorkloadConfig)
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DaskClusterWorker.
func (in *DaskClusterWorker) DeepCopy() *DaskClusterWorker {
	if in == nil {
		return nil
	}
	out := new(DaskClusterWorker)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FrameworkConfig) DeepCopyInto(out *FrameworkConfig) {
	*out = *in
	if in.Configs != nil {
		in, out := &in.Configs, &out.Configs
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FrameworkConfig.
func (in *FrameworkConfig) DeepCopy() *FrameworkConfig {
	if in == nil {
		return nil
	}
	out := new(FrameworkConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IstioConfig) DeepCopyInto(out *IstioConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IstioConfig.
func (in *IstioConfig) DeepCopy() *IstioConfig {
	if in == nil {
		return nil
	}
	out := new(IstioConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KerberosKeytabConfig) DeepCopyInto(out *KerberosKeytabConfig) {
	*out = *in
	if in.Contents != nil {
		in, out := &in.Contents, &out.Contents
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KerberosKeytabConfig.
func (in *KerberosKeytabConfig) DeepCopy() *KerberosKeytabConfig {
	if in == nil {
		return nil
	}
	out := new(KerberosKeytabConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeyTabConfig) DeepCopyInto(out *KeyTabConfig) {
	*out = *in
	if in.KeyTab != nil {
		in, out := &in.KeyTab, &out.KeyTab
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeyTabConfig.
func (in *KeyTabConfig) DeepCopy() *KeyTabConfig {
	if in == nil {
		return nil
	}
	out := new(KeyTabConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkPolicyConfig) DeepCopyInto(out *NetworkPolicyConfig) {
	*out = *in
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
	if in.ClientLabels != nil {
		in, out := &in.ClientLabels, &out.ClientLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.DashboardLabels != nil {
		in, out := &in.DashboardLabels, &out.DashboardLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkPolicyConfig.
func (in *NetworkPolicyConfig) DeepCopy() *NetworkPolicyConfig {
	if in == nil {
		return nil
	}
	out := new(NetworkPolicyConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OCIImageDefinition) DeepCopyInto(out *OCIImageDefinition) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OCIImageDefinition.
func (in *OCIImageDefinition) DeepCopy() *OCIImageDefinition {
	if in == nil {
		return nil
	}
	out := new(OCIImageDefinition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PersistentVolumeClaimTemplate) DeepCopyInto(out *PersistentVolumeClaimTemplate) {
	*out = *in
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PersistentVolumeClaimTemplate.
func (in *PersistentVolumeClaimTemplate) DeepCopy() *PersistentVolumeClaimTemplate {
	if in == nil {
		return nil
	}
	out := new(PersistentVolumeClaimTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RayCluster) DeepCopyInto(out *RayCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RayCluster.
func (in *RayCluster) DeepCopy() *RayCluster {
	if in == nil {
		return nil
	}
	out := new(RayCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RayCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RayClusterHead) DeepCopyInto(out *RayClusterHead) {
	*out = *in
	in.RayClusterNode.DeepCopyInto(&out.RayClusterNode)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RayClusterHead.
func (in *RayClusterHead) DeepCopy() *RayClusterHead {
	if in == nil {
		return nil
	}
	out := new(RayClusterHead)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RayClusterList) DeepCopyInto(out *RayClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RayCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RayClusterList.
func (in *RayClusterList) DeepCopy() *RayClusterList {
	if in == nil {
		return nil
	}
	out := new(RayClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RayClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RayClusterNetworkPolicy) DeepCopyInto(out *RayClusterNetworkPolicy) {
	*out = *in
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
	if in.ClientServerLabels != nil {
		in, out := &in.ClientServerLabels, &out.ClientServerLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.DashboardLabels != nil {
		in, out := &in.DashboardLabels, &out.DashboardLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RayClusterNetworkPolicy.
func (in *RayClusterNetworkPolicy) DeepCopy() *RayClusterNetworkPolicy {
	if in == nil {
		return nil
	}
	out := new(RayClusterNetworkPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RayClusterNode) DeepCopyInto(out *RayClusterNode) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.InitContainers != nil {
		in, out := &in.InitContainers, &out.InitContainers
		*out = make([]v1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VolumeMounts != nil {
		in, out := &in.VolumeMounts, &out.VolumeMounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VolumeClaimTemplates != nil {
		in, out := &in.VolumeClaimTemplates, &out.VolumeClaimTemplates
		*out = make([]PersistentVolumeClaimTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RayClusterNode.
func (in *RayClusterNode) DeepCopy() *RayClusterNode {
	if in == nil {
		return nil
	}
	out := new(RayClusterNode)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RayClusterSpec) DeepCopyInto(out *RayClusterSpec) {
	*out = *in
	if in.Image != nil {
		in, out := &in.Image, &out.Image
		*out = new(OCIImageDefinition)
		**out = **in
	}
	if in.ImagePullSecrets != nil {
		in, out := &in.ImagePullSecrets, &out.ImagePullSecrets
		*out = make([]v1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
	if in.Autoscaling != nil {
		in, out := &in.Autoscaling, &out.Autoscaling
		*out = new(Autoscaling)
		(*in).DeepCopyInto(*out)
	}
	in.NetworkPolicy.DeepCopyInto(&out.NetworkPolicy)
	if in.RedisShardPorts != nil {
		in, out := &in.RedisShardPorts, &out.RedisShardPorts
		*out = make([]int32, len(*in))
		copy(*out, *in)
	}
	if in.WorkerPorts != nil {
		in, out := &in.WorkerPorts, &out.WorkerPorts
		*out = make([]int32, len(*in))
		copy(*out, *in)
	}
	if in.ObjectStoreMemoryBytes != nil {
		in, out := &in.ObjectStoreMemoryBytes, &out.ObjectStoreMemoryBytes
		*out = new(int64)
		**out = **in
	}
	if in.EnableDashboard != nil {
		in, out := &in.EnableDashboard, &out.EnableDashboard
		*out = new(bool)
		**out = **in
	}
	if in.PodSecurityContext != nil {
		in, out := &in.PodSecurityContext, &out.PodSecurityContext
		*out = new(v1.PodSecurityContext)
		(*in).DeepCopyInto(*out)
	}
	if in.EnvVars != nil {
		in, out := &in.EnvVars, &out.EnvVars
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.IstioConfig = in.IstioConfig
	in.Head.DeepCopyInto(&out.Head)
	in.Worker.DeepCopyInto(&out.Worker)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RayClusterSpec.
func (in *RayClusterSpec) DeepCopy() *RayClusterSpec {
	if in == nil {
		return nil
	}
	out := new(RayClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RayClusterStatus) DeepCopyInto(out *RayClusterStatus) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RayClusterStatus.
func (in *RayClusterStatus) DeepCopy() *RayClusterStatus {
	if in == nil {
		return nil
	}
	out := new(RayClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RayClusterWorker) DeepCopyInto(out *RayClusterWorker) {
	*out = *in
	in.RayClusterNode.DeepCopyInto(&out.RayClusterNode)
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RayClusterWorker.
func (in *RayClusterWorker) DeepCopy() *RayClusterWorker {
	if in == nil {
		return nil
	}
	out := new(RayClusterWorker)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceAccountConfig) DeepCopyInto(out *ServiceAccountConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceAccountConfig.
func (in *ServiceAccountConfig) DeepCopy() *ServiceAccountConfig {
	if in == nil {
		return nil
	}
	out := new(ServiceAccountConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SparkAdditionalStorage) DeepCopyInto(out *SparkAdditionalStorage) {
	*out = *in
	if in.AccessModes != nil {
		in, out := &in.AccessModes, &out.AccessModes
		*out = make([]v1.PersistentVolumeAccessMode, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SparkAdditionalStorage.
func (in *SparkAdditionalStorage) DeepCopy() *SparkAdditionalStorage {
	if in == nil {
		return nil
	}
	out := new(SparkAdditionalStorage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SparkCluster) DeepCopyInto(out *SparkCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SparkCluster.
func (in *SparkCluster) DeepCopy() *SparkCluster {
	if in == nil {
		return nil
	}
	out := new(SparkCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SparkCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SparkClusterDriver) DeepCopyInto(out *SparkClusterDriver) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SparkClusterDriver.
func (in *SparkClusterDriver) DeepCopy() *SparkClusterDriver {
	if in == nil {
		return nil
	}
	out := new(SparkClusterDriver)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SparkClusterList) DeepCopyInto(out *SparkClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SparkCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SparkClusterList.
func (in *SparkClusterList) DeepCopy() *SparkClusterList {
	if in == nil {
		return nil
	}
	out := new(SparkClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SparkClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SparkClusterMaster) DeepCopyInto(out *SparkClusterMaster) {
	*out = *in
	in.SparkClusterNode.DeepCopyInto(&out.SparkClusterNode)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SparkClusterMaster.
func (in *SparkClusterMaster) DeepCopy() *SparkClusterMaster {
	if in == nil {
		return nil
	}
	out := new(SparkClusterMaster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SparkClusterNetworkPolicy) DeepCopyInto(out *SparkClusterNetworkPolicy) {
	*out = *in
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
	if in.ClientServerLabels != nil {
		in, out := &in.ClientServerLabels, &out.ClientServerLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.DashboardLabels != nil {
		in, out := &in.DashboardLabels, &out.DashboardLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ExternalPolicyEnabled != nil {
		in, out := &in.ExternalPolicyEnabled, &out.ExternalPolicyEnabled
		*out = new(bool)
		**out = **in
	}
	if in.ExternalPodLabels != nil {
		in, out := &in.ExternalPodLabels, &out.ExternalPodLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SparkClusterNetworkPolicy.
func (in *SparkClusterNetworkPolicy) DeepCopy() *SparkClusterNetworkPolicy {
	if in == nil {
		return nil
	}
	out := new(SparkClusterNetworkPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SparkClusterNode) DeepCopyInto(out *SparkClusterNode) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.InitContainers != nil {
		in, out := &in.InitContainers, &out.InitContainers
		*out = make([]v1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VolumeMounts != nil {
		in, out := &in.VolumeMounts, &out.VolumeMounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
	if in.AdditionalStorage != nil {
		in, out := &in.AdditionalStorage, &out.AdditionalStorage
		*out = make([]SparkAdditionalStorage, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.FrameworkConfig != nil {
		in, out := &in.FrameworkConfig, &out.FrameworkConfig
		*out = new(FrameworkConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.KeyTabConfig != nil {
		in, out := &in.KeyTabConfig, &out.KeyTabConfig
		*out = new(KeyTabConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SparkClusterNode.
func (in *SparkClusterNode) DeepCopy() *SparkClusterNode {
	if in == nil {
		return nil
	}
	out := new(SparkClusterNode)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SparkClusterSpec) DeepCopyInto(out *SparkClusterSpec) {
	*out = *in
	if in.Image != nil {
		in, out := &in.Image, &out.Image
		*out = new(OCIImageDefinition)
		**out = **in
	}
	if in.ImagePullSecrets != nil {
		in, out := &in.ImagePullSecrets, &out.ImagePullSecrets
		*out = make([]v1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
	out.IstioConfig = in.IstioConfig
	out.Driver = in.Driver
	if in.EnableDashboard != nil {
		in, out := &in.EnableDashboard, &out.EnableDashboard
		*out = new(bool)
		**out = **in
	}
	in.NetworkPolicy.DeepCopyInto(&out.NetworkPolicy)
	if in.PodSecurityContext != nil {
		in, out := &in.PodSecurityContext, &out.PodSecurityContext
		*out = new(v1.PodSecurityContext)
		(*in).DeepCopyInto(*out)
	}
	if in.EnvVars != nil {
		in, out := &in.EnvVars, &out.EnvVars
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Master.DeepCopyInto(&out.Master)
	in.Worker.DeepCopyInto(&out.Worker)
	if in.Autoscaling != nil {
		in, out := &in.Autoscaling, &out.Autoscaling
		*out = new(Autoscaling)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SparkClusterSpec.
func (in *SparkClusterSpec) DeepCopy() *SparkClusterSpec {
	if in == nil {
		return nil
	}
	out := new(SparkClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SparkClusterStatus) DeepCopyInto(out *SparkClusterStatus) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SparkClusterStatus.
func (in *SparkClusterStatus) DeepCopy() *SparkClusterStatus {
	if in == nil {
		return nil
	}
	out := new(SparkClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SparkClusterWorker) DeepCopyInto(out *SparkClusterWorker) {
	*out = *in
	in.SparkClusterNode.DeepCopyInto(&out.SparkClusterNode)
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SparkClusterWorker.
func (in *SparkClusterWorker) DeepCopy() *SparkClusterWorker {
	if in == nil {
		return nil
	}
	out := new(SparkClusterWorker)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkloadConfig) DeepCopyInto(out *WorkloadConfig) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.InitContainers != nil {
		in, out := &in.InitContainers, &out.InitContainers
		*out = make([]v1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VolumeMounts != nil {
		in, out := &in.VolumeMounts, &out.VolumeMounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VolumeClaimTemplates != nil {
		in, out := &in.VolumeClaimTemplates, &out.VolumeClaimTemplates
		*out = make([]PersistentVolumeClaimTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkloadConfig.
func (in *WorkloadConfig) DeepCopy() *WorkloadConfig {
	if in == nil {
		return nil
	}
	out := new(WorkloadConfig)
	in.DeepCopyInto(out)
	return out
}
