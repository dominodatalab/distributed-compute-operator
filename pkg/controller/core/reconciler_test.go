package core

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"

	"github.com/stretchr/testify/assert"

	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

type dummyComponent struct {
	kind client.Object
}

func newOwnedComponent(kind client.Object) Component {
	return &dummyComponent{kind: kind}
}

func (c dummyComponent) Reconcile(ctx *Context) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func (c dummyComponent) Kind() client.Object { // From OwnedComponent
	return c.kind
}

func (c dummyComponent) Finalize(*Context) (ctrl.Result, bool, error) { // From FinalizerComponent
	return ctrl.Result{}, false, nil
}

var mgr ctrl.Manager

func TestMain(m *testing.M) {
	var env = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := env.Start()
	if err != nil {
		fmt.Printf("Cannot start the environment: %v", err)
		os.Exit(1)
	}

	err = dcv1alpha1.AddToScheme(scheme.Scheme)
	if err != nil {
		// Probably OK here, but some tests will fail later
		fmt.Printf("Schema failed: %v", err)
	}

	code := 0

	mgr, err = ctrl.NewManager(cfg, ctrl.Options{Scheme: scheme.Scheme})
	if err != nil {
		fmt.Printf("Manager initialization failed: %v", err)
		code = 1
	} else {
		code = m.Run()
	}

	err = env.Stop()
	if err != nil {
		// This is OK
		fmt.Printf("Cannot stop the environment: %v", err)
	}

	os.Exit(code)
}

func TestNewReconciler(t *testing.T) {
	reconciler := NewReconciler(mgr)
	assert.Equal(t, mgr, reconciler.mgr)
	assert.Equal(t, mgr.GetClient(), reconciler.client)
	assert.Empty(t, reconciler.components)
	assert.False(t, reconciler.webhooksEnabled)
}

func TestFor(t *testing.T) {
	reconciler := NewReconciler(mgr)
	apiType := &dcv1alpha1.MPICluster{}
	reconciler.For(apiType)
	assert.Equal(t, apiType, reconciler.apiType)
}

func TestWithWebHooks(t *testing.T) {
	reconciler := NewReconciler(mgr)
	reconciler.WithWebhooks()
	assert.True(t, reconciler.webhooksEnabled)
}

func TestComponents(t *testing.T) {
	reconciler := NewReconciler(mgr)
	reconciler.Component("ConfigMapComponent", newOwnedComponent(&corev1.ConfigMap{}))
	reconciler.Component("PodComponent", newOwnedComponent(&corev1.Pod{}))
	assert.Equal(t, 2, len(reconciler.components))
}

func TestComplete(t *testing.T) {
	reconciler := NewReconciler(mgr).
		For(&dcv1alpha1.MPICluster{}).
		Component("ConfigMapComponent", newOwnedComponent(&corev1.ConfigMap{})).
		Component("PodComponent", newOwnedComponent(&corev1.Pod{}))

	assert.Nil(t, reconciler.Complete())

	assert.Equal(t, "mpicluster", reconciler.name)
	assert.NotNil(t, reconciler.log)
}

func TestCompleteFailsWithoutComponents(t *testing.T) {
	reconciler := NewReconciler(mgr).
		For(&dcv1alpha1.MPICluster{})

	assert.NotNil(t, reconciler.Complete())
}

func TestCompleteFailsWithoutAPIType(t *testing.T) {
	reconciler := NewReconciler(mgr).
		Component("ConfigMapComponent", newOwnedComponent(&corev1.ConfigMap{})).
		Component("PodComponent", newOwnedComponent(&corev1.Pod{}))

	assert.NotNil(t, reconciler.Complete())
}
