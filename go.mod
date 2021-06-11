module github.com/dominodatalab/distributed-compute-operator

go 1.16

require (
	github.com/banzaicloud/k8s-objectmatcher v1.5.1
	github.com/blang/semver v3.5.1+incompatible
	github.com/docker/distribution v2.7.1+incompatible
	github.com/go-logr/logr v0.4.0
	github.com/hashicorp/go-retryablehttp v0.6.8
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/prometheus/common v0.26.0
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.7.0
	istio.io/api v0.0.0-20210318170531-e6e017e575c5
	istio.io/client-go v1.9.2
	k8s.io/api v0.21.1
	k8s.io/apiextensions-apiserver v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/apiserver v0.21.1
	k8s.io/client-go v0.21.1
	k8s.io/utils v0.0.0-20210527160623-6fdb442a123b
	sigs.k8s.io/controller-runtime v0.9.0
	sigs.k8s.io/yaml v1.2.0
)
