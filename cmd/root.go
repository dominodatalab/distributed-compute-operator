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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// NOTE: https://github.com/spf13/cobra/issues/587
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}