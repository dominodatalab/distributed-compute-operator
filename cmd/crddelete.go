package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/dominodatalab/distributed-compute-operator/pkg/crd"
)

var crdDeleteCmd = &cobra.Command{
	Use:   "crd-delete",
	Short: "Delete custom resource definitions from a cluster",
	Long: `Delete all "distributed-compute.dominodatalab.com" CRDs from a cluster.

Any running distributed compute resources will be decommissioned when this 
operation runs (i.e. your deployments will be deleted immediately). This will
only attempt to remove definitions that are already present in Kubernetes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return crd.Delete(context.Background(), istioEnabled)
	},
}

func init() {
	rootCmd.AddCommand(crdDeleteCmd)
}
