package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "distributed-compute-operator",
	Short: "Kubernetes operator that manages parallel computing clusters.",
	Long:  `Kubernetes operator that manages parallel computing clusters.`,
}

// Execute launches the command line tool.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func addIstioFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("istio-enabled", "i", false, "Enable support for Istio sidecar container")
}

func processIstioFlag(op func(enabled bool) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		istioEnabled, err := cmd.Flags().GetBool("istio-enabled")
		if err != nil {
			return err
		}

		return op(istioEnabled)
	}
}

func init() {
	// NOTE: required until https://github.com/spf13/cobra/issues/587
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}
