package core

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/dominodatalab/distributed-compute-operator/pkg/logging"
)

var getGvk = apiutil.GVKForObject

const skipReconcileAnnotation = "controller.dominodatalab.com/skip-reconcile"

type reconcilerComponent struct {
	name string
	comp Component

	finalizer     FinalizerComponent
	finalizerName string
}

type Reconciler struct {
	name              string
	mgr               ctrl.Manager
	controllerBuilder *ctrl.Builder
	apiType           client.Object
	client            client.Client
	log               logr.Logger
	webhooksEnabled   bool
	finalizerBaseName string

	patcher    *Patch
	recorder   record.EventRecorder
	controller controller.Controller
	components []*reconcilerComponent
}

func NewReconciler(mgr ctrl.Manager) *Reconciler {
	return &Reconciler{
		mgr:               mgr,
		client:            mgr.GetClient(),
		components:        []*reconcilerComponent{},
		controllerBuilder: builder.ControllerManagedBy(mgr),
	}
}

func (r *Reconciler) For(apiType client.Object, opts ...builder.ForOption) *Reconciler {
	r.apiType = apiType
	r.controllerBuilder = r.controllerBuilder.For(apiType, opts...)

	return r
}

func (r *Reconciler) Component(name string, comp Component, opts ...builder.OwnsOption) *Reconciler {
	rc := &reconcilerComponent{name: name, comp: comp}
	r.controllerBuilder.Owns(comp.Kind(), opts...)

	if finalizer, ok := comp.(FinalizerComponent); ok {
		rc.finalizer = finalizer
	}
	r.components = append(r.components, rc)

	return r
}

func (r *Reconciler) WithWebhooks() *Reconciler {
	r.webhooksEnabled = true
	return r
}

func (r *Reconciler) Build() (controller.Controller, error) {
	name, err := r.getControllerName()
	if err != nil {
		return nil, fmt.Errorf("cannot compute controller name: %w", err)
	}
	r.name = name
	r.log = logging.New(ctrl.Log.WithName("controllers").WithName(name))
	r.recorder = r.mgr.GetEventRecorderFor(fmt.Sprintf("%s-%s", r.name, "controller"))

	if r.finalizerBaseName == "" {
		gvk, err := getGvk(r.apiType, r.mgr.GetScheme())
		if err != nil {
			return nil, fmt.Errorf("cannot get GVK for object %#v: %w", r.apiType, err)
		}
		r.finalizerBaseName = fmt.Sprintf("%s.%s/", name, gvk.Group)
	}

	if r.patcher == nil {
		gvk, err := getGvk(r.apiType, r.mgr.GetScheme())
		if err != nil {
			return nil, fmt.Errorf("cannot get GVK for object %#v: %w", r.apiType, err)
		}
		r.patcher = NewPatch(gvk)
	}

	components := map[string]Component{}
	for _, rc := range r.components {
		orig, ok := components[rc.name]
		if ok {
			return nil, fmt.Errorf("duplicate component found using name %s: %#v %#v", rc.name, orig, rc.comp)
		}
		rc.finalizerName = path.Join(r.finalizerBaseName, rc.name)

		components[rc.name] = rc.comp
	}

	r.controller, err = r.controllerBuilder.Build(r)
	if err != nil {
		return nil, fmt.Errorf("unable to build controller: %w", err)
	}

	// setup webhooks
	if r.webhooksEnabled {
		err := ctrl.NewWebhookManagedBy(r.mgr).For(r.apiType).Complete()
		if err != nil {
			return nil, fmt.Errorf("unable to build webhook: %w", err)
		}
	}

	return r.controller, nil
}

func (r *Reconciler) Complete() error {
	_, err := r.Build()
	return err
}

func (r *Reconciler) Reconcile(rootCtx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.log.WithValues(r.name, req.NamespacedName)

	// fetch event api object
	log.Info("Starting reconcile")

	obj := r.apiType.DeepCopyObject().(client.Object)
	if err := r.client.Get(rootCtx, req.NamespacedName, obj); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Aborting reconcile, object not found (assuming it was deleted)")
			return ctrl.Result{}, nil
		}

		log.Error(err, "Failed to fetch reconcile object")
		return ctrl.Result{}, err
	}
	cleanObj := obj.DeepCopyObject().(client.Object)

	// skip reconcile when annotated
	skip, ok := obj.GetAnnotations()[skipReconcileAnnotation]
	if ok && skip == "true" {
		log.Info("Skipping reconcile due to annotation")
		return ctrl.Result{}, nil
	}

	// build context for components
	compLog := log.WithName("components")
	ctx := &Context{
		Context:  rootCtx,
		Object:   obj,
		Client:   r.client,
		Patch:    r.patcher,
		Scheme:   r.mgr.GetScheme(),
		Recorder: r.recorder,
	}

	// reconcile components
	var finalRes ctrl.Result
	var errs []error
	for _, rc := range r.components {
		res := ctrl.Result{}
		var err error

		ctx.Log = compLog.WithName(rc.name)

		if ctx.Object.GetDeletionTimestamp().IsZero() {
			log.Info("Reconciling component", "component", rc.name)
			res, err = rc.comp.Reconcile(ctx)

			if rc.finalizer != nil {
				log.Info("Registering finalizer", "component", rc.name)
				controllerutil.AddFinalizer(ctx.Object, rc.finalizerName)
			}
		} else if rc.finalizer != nil && controllerutil.ContainsFinalizer(ctx.Object, rc.finalizerName) {
			log.Info("Finalizing component", "component", rc.name)

			var done bool
			res, done, err = rc.finalizer.Finalize(ctx)
			if done {
				log.Info("Removing finalizer", "component", rc.name)
				controllerutil.RemoveFinalizer(ctx.Object, rc.finalizerName)
			}
		}

		if res.Requeue {
			finalRes.Requeue = true
		}
		if res.RequeueAfter != 0 && (finalRes.RequeueAfter == 0 || finalRes.RequeueAfter > res.RequeueAfter) {
			finalRes.RequeueAfter = res.RequeueAfter
		}
		if err != nil {
			log.Error(err, "Component reconciliation failed", "component", rc.name)
			errs = append(errs, err)
		}
	}

	// patch metadata when changes occur
	currentMeta := r.apiType.DeepCopyObject().(client.Object)
	currentMeta.SetName(ctx.Object.GetName())
	currentMeta.SetNamespace(ctx.Object.GetNamespace())
	currentMeta.SetLabels(ctx.Object.GetLabels())
	currentMeta.SetAnnotations(ctx.Object.GetAnnotations())
	currentMeta.SetFinalizers(ctx.Object.GetFinalizers())

	cleanMeta := r.apiType.DeepCopyObject().(client.Object)
	cleanMeta.SetName(cleanObj.GetName())
	cleanMeta.SetNamespace(cleanObj.GetNamespace())
	cleanMeta.SetLabels(cleanObj.GetLabels())
	cleanMeta.SetAnnotations(cleanObj.GetAnnotations())
	cleanMeta.SetFinalizers(cleanObj.GetFinalizers())

	patch := client.MergeFrom(cleanMeta)
	json, _ := patch.Data(currentMeta)

	log.V(1).Info("Patching metadata", "type", patch.Type(), "patch", string(json))
	if err := r.client.Patch(ctx, currentMeta, patch, &client.PatchOptions{FieldManager: r.name}); err != nil {
		return ctrl.Result{}, fmt.Errorf("error patching metadata: %w", err)
	}

	// condense all error messages into one
	var err error
	if len(errs) == 1 {
		err = errs[0]
	} else if len(errs) > 1 {
		sb := strings.Builder{}
		for _, e := range errs {
			sb.WriteString(fmt.Sprintf(" %s\n", e.Error()))
		}
		err = errors.New(sb.String())
	}

	return finalRes, err
}

func (r *Reconciler) getControllerName() (string, error) {
	if r.name != "" {
		return r.name, nil
	}

	gvk, err := getGvk(r.apiType, r.mgr.GetScheme())
	if err != nil {
		return "", err
	}

	return strings.ToLower(gvk.Kind), nil
}
