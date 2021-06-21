package core

import (
	"path"

	"github.com/banzaicloud/k8s-objectmatcher/patch"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var defaultCalculateOpts = []patch.CalculateOption{
	patch.IgnoreStatusFields(),
	patch.IgnoreVolumeClaimTemplateTypeMetaAndStatus(),
}

type Patch struct {
	Annotator     *patch.Annotator
	Maker         *patch.PatchMaker
	CalculateOpts []patch.CalculateOption
}

func NewPatch(gvk schema.GroupVersionKind) *Patch {
	a := patch.NewAnnotator(path.Join(gvk.Group, "last-applied"))
	m := patch.NewPatchMaker(a)

	return &Patch{
		Annotator:     a,
		Maker:         m,
		CalculateOpts: defaultCalculateOpts,
	}
}
