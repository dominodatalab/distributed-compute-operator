package cmd

import (
	"flag"

	"github.com/dominodatalab/distributed-compute-operator/controllers"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/dominodatalab/distributed-compute-operator/pkg/manager"
)

const WebhookPort = 9443

var (
	namespaces           []string
	probeAddr            string
	metricsAddr          string
	webhookPort          int
	enableLeaderElection bool
	zapOpts              = zap.Options{}
	mpiInitImage         string
	mpiSyncImage         string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the controller manager",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &controllers.Config{
			Namespaces:           namespaces,
			MetricsAddr:          metricsAddr,
			HealthProbeAddr:      probeAddr,
			WebhookServerPort:    webhookPort,
			EnableLeaderElection: enableLeaderElection,
			IstioEnabled:         istioEnabled,
			ZapOptions:           zapOpts,
			MPIInitImage:         mpiInitImage,
			MPISyncImage:         mpiSyncImage,
		}

		return manager.Start(cfg)
	},
}

func init() {
	startCmd.Flags().SortFlags = false

	fs := new(flag.FlagSet)
	zapOpts.BindFlags(fs)

	startCmd.Flags().AddGoFlagSet(fs)
	startCmd.Flags().StringSliceVar(&namespaces, "namespaces", nil,
		"Only reconcile resources in these namespaces")
	startCmd.Flags().IntVar(&webhookPort, "webhook-server-port", WebhookPort,
		"Webhook server will bind to this port")
	startCmd.Flags().StringVar(&metricsAddr, "metrics-bind-address", ":8080",
		"Metrics endpoint will bind to this address")
	startCmd.Flags().StringVar(&probeAddr, "health-probe-bind-address", ":8081",
		"Health probe endpoint will bind to this address")
	startCmd.Flags().BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election to ensure there is only one active controller manager")
	startCmd.Flags().StringVar(&mpiInitImage, "mpi-init-image", "",
		"Image for MPI worker init container")
	startCmd.Flags().StringVar(&mpiSyncImage, "mpi-sync-image", "",
		"Image for MPI worker sync container")

	rootCmd.AddCommand(startCmd)
}
