package cmd

import (
	"flag"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/dominodatalab/distributed-compute-operator/pkg/manager"
)

var (
	namespace            string
	probeAddr            string
	metricsAddr          string
	webhookPort          int
	enableLeaderElection bool

	zapOpts = zap.Options{}
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the controller manager",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &manager.Config{
			Namespace:            namespace,
			MetricsAddr:          metricsAddr,
			HealthProbeAddr:      probeAddr,
			WebhookServerPort:    webhookPort,
			EnableLeaderElection: enableLeaderElection,
			IstioEnabled:         istioEnabled,
			ZapOptions:           zapOpts,
		}

		return manager.Start(cfg)
	},
}

func init() {
	startCmd.Flags().SortFlags = false

	fs := new(flag.FlagSet)
	zapOpts.BindFlags(fs)

	startCmd.Flags().AddGoFlagSet(fs)
	startCmd.Flags().StringVar(&namespace, "namespace", "default", "Reconcile clusters resources in this namespace")
	startCmd.Flags().IntVar(&webhookPort, "webhook-server-port", 9443, "Webhook server will bind to this port")
	startCmd.Flags().StringVar(&metricsAddr, "metrics-bind-address", ":8080",
		"Metrics endpoint will bind to this address")
	startCmd.Flags().StringVar(&probeAddr, "health-probe-bind-address", ":8081",
		"Health probe endpoint will bind to this address")
	startCmd.Flags().BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election to ensure there is only one active controller manager")

	rootCmd.AddCommand(startCmd)
}
