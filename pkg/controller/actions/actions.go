package actions

import (
	"k8s.io/apimachinery/pkg/runtime/schema"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

func CreateOrUpdateOwnedResource(ctx *core.Context, owner metav1.Object, controlled client.Object) error {

	err := ctrl.SetControllerReference(owner, controlled, ctx.Scheme)
	if err != nil {
		return err
	}

	gvk, err := getObjectKind(ctx, controlled)
	if err != nil {
		return err
	}

	found := controlled.DeepCopyObject().(client.Object)

	err = ctx.Client.Get(ctx, client.ObjectKeyFromObject(controlled), found)
	if err != nil {
		if apierrors.IsNotFound(err) {
			ctx.Log.V(1).Info("Creating controlled object", "gvk", gvk, "object", controlled)
			return createOwnedResource(ctx, controlled)
		}
		return err
	}

	patchResult, err := ctx.Patch.Maker.Calculate(found, controlled, ctx.Patch.CalculateOpts...)
	if err != nil {
		return err
	}
	if patchResult.IsEmpty() {
		return nil
	}

	ctx.Log.V(1).Info("Applying patch to object", "gvk", gvk, "object", controlled, "patch",
		string(patchResult.Patch))
	err = ctx.Patch.Annotator.SetLastAppliedAnnotation(controlled)
	if err != nil {
		return err
	}

	controlled.SetResourceVersion(found.GetResourceVersion())

	// ensure we do not modify "generated" values for certain resources
	switch modified := controlled.(type) {
	case *corev1.Service:
		current := found.(*corev1.Service)
		modified.Spec.ClusterIP = current.Spec.ClusterIP
	case *batchv1.Job:
		current := found.(*batchv1.Job)
		modified.Spec.Selector = current.Spec.Selector
	}

	ctx.Log.V(1).Info("Updating controlled object", "gvk", gvk, "object", controlled)
	return ctx.Client.Update(ctx, controlled)
}

func DeleteIfExists(ctx *core.Context, objs ...client.Object) error {
	for _, obj := range objs {
		if err := ctx.Client.Get(ctx, client.ObjectKeyFromObject(obj), obj); err != nil {
			if apierrors.IsNotFound(err) {
				continue
			}

			return err
		}

		ctx.Log.V(1).Info("Deleting controlled object", "object", obj)
		if err := ctx.Client.Delete(ctx, obj); err != nil {
			return err
		}
	}

	return nil
}

func DeleteStorage(ctx *core.Context, opts []client.ListOption) error {
	pvcList := &corev1.PersistentVolumeClaimList{}
	listOpts := (&client.ListOptions{}).ApplyOptions(opts)

	ctx.Log.Info("Querying for persistent volume claims", "namespace", listOpts.Namespace, "labels",
		listOpts.LabelSelector.String())
	if err := ctx.Client.List(ctx, pvcList, opts...); err != nil {
		ctx.Log.Error(err, "Cannot list persistent volume claims")
		return err
	}

	for idx := range pvcList.Items {
		pvc := &pvcList.Items[idx]
		key := client.ObjectKeyFromObject(pvc)

		ctx.Log.Info("Deleting persistent volume claim", "claim", key)
		if err := ctx.Client.Delete(ctx, pvc); err != nil {
			ctx.Log.Error(err, "Cannot delete persistent volume claim", "claim", key)
			return err
		}
	}

	return nil
}

func getObjectKind(ctx *core.Context, controlled client.Object) (schema.GroupVersionKind, error) {
	gvks, _, err := ctx.Scheme.ObjectKinds(controlled)
	if err != nil {
		return schema.GroupVersionKind{}, err
	}
	return gvks[0], nil
}

func createOwnedResource(ctx *core.Context, controlled client.Object) error {
	err := ctx.Patch.Annotator.SetLastAppliedAnnotation(controlled)
	if err != nil {
		return err
	}
	return ctx.Client.Create(ctx, controlled)
}
