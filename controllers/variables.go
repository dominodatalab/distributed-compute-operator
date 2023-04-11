package controllers

import (
	"path"

	"github.com/banzaicloud/k8s-objectmatcher/patch"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

var (
	// DistributedComputeFinalizer is the custom identifier used to mark
	// controller-managed resources that require pre-delete hook logic.
	DistributedComputeFinalizer = path.Join(dcv1alpha1.GroupVersion.Group, "finalizer")

	// PatchAnnotator applies state annotations to owned components.
	PatchAnnotator = patch.NewAnnotator(path.Join(dcv1alpha1.GroupVersion.Group, "last-applied"))
	// PatchMaker calculates changes to state annotations on owned components.
	PatchMaker = patch.NewPatchMaker(PatchAnnotator, &patch.K8sStrategicMergePatcher{}, &patch.BaseJSONMergePatcher{})
	// PatchCalculateOpts define the exclusion rules used when calculating the
	// difference between two k8s resources.
	PatchCalculateOpts = []patch.CalculateOption{
		patch.IgnoreStatusFields(),
		patch.IgnoreVolumeClaimTemplateTypeMetaAndStatus(),
	}
)
