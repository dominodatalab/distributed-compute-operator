package crd

import (
	"context"
	"encoding/json"

	apixv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apixv1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/yaml"

	"github.com/dominodatalab/distributed-compute-operator/config/crd"
)

var log = zap.New()

// Apply will create or update all project CRDs inside a Kubernetes cluster.
// The latest available version of the CRD will be used to perform this operation.
func Apply(ctx context.Context) error {
	apply := func(client apixv1client.CustomResourceDefinitionInterface, crd *apixv1.CustomResourceDefinition) error {
		found, err := client.Get(ctx, crd.Name, metav1.GetOptions{})

		switch {
		case apierrors.IsNotFound(err):
			log.Info("Creating CRD", "Name", crd.Name)
			_, err = client.Create(ctx, crd, metav1.CreateOptions{})
		case err == nil:
			log.Info("Updating CRD", "Name", crd.Name)
			crd.SetResourceVersion(found.ResourceVersion)
			_, err = client.Update(ctx, crd, metav1.UpdateOptions{})
		}

		return err
	}

	return processCRDs(apply)
}

// Delete will remove all project CRDs from a Kubernetes cluster.
func Delete(ctx context.Context) error {
	deleteFn := func(client apixv1client.CustomResourceDefinitionInterface, crd *apixv1.CustomResourceDefinition) error {
		log.Info("Deleting CRD", "Name", crd.Name)

		err := client.Delete(ctx, crd.Name, metav1.DeleteOptions{})
		if apierrors.IsNotFound(err) {
			log.Info("CRD not found, ignoring", "Name", crd.Name)
			return nil
		}

		return err
	}

	return processCRDs(deleteFn)
}

// processCRDs loads all available CRDs and uses a processor func to act upon them.
func processCRDs(processor func(client apixv1client.CustomResourceDefinitionInterface, crd *apixv1.CustomResourceDefinition) error) error {
	client, err := getCRDClient()
	if err != nil {
		return err
	}

	log.Info("Loading CRDs")
	definitions, err := crd.ReadAll()
	if err != nil {
		return err
	}

	for _, def := range definitions {
		crd, err := loadCRD(def)
		if err != nil {
			return err
		}

		if err := processor(client, crd); err != nil {
			return err
		}
	}

	return nil
}

// loadCRD converts YAML bytes to JSON and loads that into a CustomResourceDefinition object.
func loadCRD(bs []byte) (*apixv1.CustomResourceDefinition, error) {
	bs, err := yaml.YAMLToJSON(bs)
	if err != nil {
		return nil, err
	}

	crd := new(apixv1.CustomResourceDefinition)
	if err := json.Unmarshal(bs, crd); err != nil {
		return nil, err
	}

	return crd, nil
}

// getCRDClient returns a client configured to work with custom resource definitions.
func getCRDClient() (apixv1client.CustomResourceDefinitionInterface, error) {
	log.Info("Initializing Kubernetes CRD client")

	config, err := loadKubernetesConfig()
	if err != nil {
		return nil, err
	}

	client, err := apixv1client.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client.CustomResourceDefinitions(), nil
}

// loadKubernetesConfig returns the canonical kubernetes config for use with a client.
// Out-of-cluster loading is attempted first, followed by in-cluster when that fails.
func loadKubernetesConfig() (*rest.Config, error) {
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	if cfg, err := kubeconfig.ClientConfig(); err == nil {
		return cfg, nil
	}

	return rest.InClusterConfig()
}
