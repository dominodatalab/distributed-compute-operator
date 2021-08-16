package manager

import "sigs.k8s.io/controller-runtime/pkg/log/zap"

// Config options for the controller manager.
type Config struct {
	MetricsAddr          string
	HealthProbeAddr      string
	WebhookServerPort    int
	EnableLeaderElection bool
	IstioEnabled         bool
	ZapOptions           zap.Options
}
