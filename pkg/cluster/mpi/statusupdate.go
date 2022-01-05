package mpi

import (
	"fmt"
	"reflect"
	"sort"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

const (
	ReadyStatus   = "Ready"
	PendingStatus = "Pending"
)

func StatusUpdate() core.Component {
	return &statusUpdateComponent{}
}

type statusUpdateComponent struct{}

func (c statusUpdateComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	cr := objToMPICluster(ctx.Object)

	var modified bool

	// update image reference
	image, err := util.ParseImageDefinition(cr.Spec.Image)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("cannot build cluster image: %w", err)
	}
	if cr.Status.Image != image {
		cr.Status.Image = image
		modified = true
	}

	// update node list
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(cr.Namespace),
		client.MatchingLabels(meta.StandardLabels(cr)),
	}
	if lErr := ctx.Client.List(ctx, podList, listOpts...); lErr != nil {
		return ctrl.Result{}, fmt.Errorf("cannot list cluster pods: %w", lErr)
	}

	var podNames []string
	for _, pod := range podList.Items {
		podNames = append(podNames, pod.Name)
	}
	sort.Strings(podNames)

	if !reflect.DeepEqual(podNames, cr.Status.Nodes) {
		cr.Status.Nodes = podNames
		modified = true
	}

	actualPodCnt, err := getActivePodCnt(ctx, cr)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("cannot obtain active pod count: %w", err)
	}
	expectedPodCnt := int(*cr.Spec.Worker.Replicas)

	status := cr.Status.ClusterStatus
	if actualPodCnt >= expectedPodCnt {
		status = ReadyStatus
	} else {
		status = PendingStatus
	}

	if cr.Status.ClusterStatus != status {
		modified = true
		cr.Status.ClusterStatus = status
		if status == ReadyStatus {
			tt := metav1.Now()
			cr.Status.StartTime = &tt
		} else {
			cr.Status.StartTime = nil
		}
	}

	if modified {
		err = ctx.Client.Status().Update(ctx, cr)
	}

	// TODO: scale down workers after job is "Complete" or "Failed"
	// 	this component focuses on updating the "status" so the downscaling
	// 	logic probably belongs elsewhere

	return ctrl.Result{}, err
}

func getActivePodCnt(ctx *core.Context, cr *dcv1alpha1.MPICluster) (int, error) {
	var ep corev1.Endpoints
	objKey := client.ObjectKey{
		Name:      serviceName(cr),
		Namespace: cr.Namespace,
	}

	if err := ctx.Client.Get(ctx, objKey, &ep); err != nil && !apierrors.IsNotFound(err) {
		return 0, fmt.Errorf("cannot fetch endpoints: %w", err)
	}

	if len(ep.Subsets) == 1 {
		return len(ep.Subsets[0].Addresses), nil
	}
	return 0, nil
}
