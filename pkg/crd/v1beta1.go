package crd

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/blang/semver"
	apixv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/yaml"

	"github.com/dominodatalab/distributed-compute-operator/config/crd"
)

// nolint:lll
type v1Beta1CRDProcessor func(context.Context, apixv1beta1client.CustomResourceDefinitionInterface, *apixv1beta1.CustomResourceDefinition) error

var (
	v1VersionConstraint = semver.MustParse("1.16.0")

	// nolint:dupl,lll
	applyV1Beta1Fn = func(ctx context.Context, client apixv1beta1client.CustomResourceDefinitionInterface, crd *apixv1beta1.CustomResourceDefinition) error {
		log.Info("Fetching CRD", "Name", crd.Name)
		found, err := client.Get(ctx, crd.Name, metav1.GetOptions{})

		if apierrors.IsNotFound(err) {
			log.Info("CRD not found, creating", "Name", crd.Name)
			_, err = client.Create(ctx, crd, metav1.CreateOptions{})
		} else if err == nil {
			log.Info("CRD found, updating", "Name", crd.Name)
			crd.SetResourceVersion(found.ResourceVersion)
			_, err = client.Update(ctx, crd, metav1.UpdateOptions{})
		}

		return err
	}

	// nolint:lll
	deleteV1Beta1Fn = func(ctx context.Context, client apixv1beta1client.CustomResourceDefinitionInterface, crd *apixv1beta1.CustomResourceDefinition) error {
		log.Info("Deleting CRD", "Name", crd.Name)
		err := client.Delete(ctx, crd.Name, metav1.DeleteOptions{})

		if apierrors.IsNotFound(err) {
			log.Info("CRD not found, ignoring", "Name", crd.Name)
			return nil
		}

		return err
	}
)

func isV1CRDAPIAvailable() (useV1 bool, err error) {
	config, err := loadKubernetesConfig()
	if err != nil {
		return
	}

	client, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return
	}

	info, err := client.ServerVersion()
	if err != nil {
		return
	}

	version, err := semver.Parse(strings.TrimPrefix(info.GitVersion, "v"))
	if err != nil {
		return
	}

	return version.GTE(v1VersionConstraint), nil
}

func processV1Beta1CRDs(ctx context.Context, processor v1Beta1CRDProcessor, definitions []crd.Definition) error {
	client, err := getV1Beta1CRDClient()
	if err != nil {
		return err
	}

	for _, def := range definitions {
		customResourceDefinition, err := loadV1Beta1CRD(def.Contents)
		if err != nil {
			return err
		}

		if err := processor(ctx, client, customResourceDefinition); err != nil {
			return err
		}
	}

	return nil
}

func loadV1Beta1CRD(bs []byte) (*apixv1beta1.CustomResourceDefinition, error) {
	bs, err := yaml.YAMLToJSON(bs)
	if err != nil {
		return nil, err
	}

	resource := new(apixv1beta1.CustomResourceDefinition)
	if err := json.Unmarshal(bs, resource); err != nil {
		return nil, err
	}

	return resource, nil
}

func getV1Beta1CRDClient() (apixv1beta1client.CustomResourceDefinitionInterface, error) {
	log.Info("Initializing Kubernetes V1Beta1 CRD client")

	config, err := loadKubernetesConfig()
	if err != nil {
		return nil, err
	}

	client, err := apixv1beta1client.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client.CustomResourceDefinitions(), nil
}
