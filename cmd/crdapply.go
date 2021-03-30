package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/dominodatalab/distributed-compute-operator/pkg/crd"
)

var crdApplyCmd = &cobra.Command{
	Use:   "crd-apply",
	Short: "Apply custom resource definitions to a cluster",
	Long: `Apply all "distributed-compute.dominodatalab.com" CRDs to a cluster.

Apply Rules:
  - When a definition is is missing, it will be created
  - If a definition is already present, then it will be updated
  - Updating definitions that have not changed results in a no-op`,
	RunE: processIstioFlag(func(enabled bool) error {
		return crd.Apply(context.Background(), enabled)
	}),
}

func init() {
	addIstioFlag(crdApplyCmd)
	rootCmd.AddCommand(crdApplyCmd)
}
