module github.com/dominodatalab/distributed-compute-operator

go 1.16

require (
	github.com/banzaicloud/k8s-objectmatcher v1.5.1
	github.com/docker/distribution v2.7.1+incompatible
	github.com/go-logr/logr v0.3.0
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.5
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.7.0
	k8s.io/api v0.19.2
	k8s.io/apiextensions-apiserver v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/utils v0.0.0-20200912215256-4140de9c8800
	sigs.k8s.io/controller-runtime v0.7.0
	sigs.k8s.io/yaml v1.2.0
)
