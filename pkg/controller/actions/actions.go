package actions

import (
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

func CreateOwnedResource(ctx *core.Context, owner metav1.Object, controlled client.Object) error {
	_, _, err := createOwnedResource(ctx, owner, controlled)
	return err
}

func CreateOrUpdateOwnedResource(ctx *core.Context, owner metav1.Object, controlled client.Object) error {
	found, gvk, err := createOwnedResource(ctx, owner, controlled)
	if err != nil {
		return err
	}

	patchResult, err := ctx.Patch.Maker.Calculate(found, controlled, ctx.Patch.CalculateOpts...)
	if err != nil {
		return err
	}
	if patchResult.IsEmpty() {
		return nil
	}

	ctx.Log.V(1).Info("Applying patch to object", "gvk", gvk, "object", controlled, "patch", string(patchResult.Patch))
	if err = ctx.Patch.Annotator.SetLastAppliedAnnotation(controlled); err != nil {
		return err
	}

	controlled.SetResourceVersion(found.GetResourceVersion())
	if modified, ok := controlled.(*corev1.Service); ok {
		current := found.(*corev1.Service)
		modified.Spec.ClusterIP = current.Spec.ClusterIP
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

	ctx.Log.Info("Querying for persistent volume claims", "namespace", listOpts.Namespace, "labels", listOpts.LabelSelector.String())
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

func createOwnedResource(
	ctx *core.Context,
	owner metav1.Object,
	controlled client.Object,
) (found client.Object, gvk schema.GroupVersionKind, err error) {
	if err = ctrl.SetControllerReference(owner, controlled, ctx.Scheme); err != nil {
		return
	}

	var gvks []schema.GroupVersionKind
	if gvks, _, err = ctx.Scheme.ObjectKinds(controlled); err != nil {
		return
	}
	gvk = gvks[0]

	found = controlled.DeepCopyObject().(client.Object)
	if err = ctx.Client.Get(ctx, client.ObjectKeyFromObject(controlled), found); err != nil {
		if !apierrors.IsNotFound(err) {
			return
		}

		if err = ctx.Patch.Annotator.SetLastAppliedAnnotation(controlled); err != nil {
			return
		}

		ctx.Log.V(1).Info("Creating controlled object", "gvk", gvk, "object", controlled)
		err = ctx.Client.Create(ctx, controlled)
	}

	return
}
