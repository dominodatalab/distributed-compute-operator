package manager

import "sigs.k8s.io/controller-runtime/pkg/log/zap"

type Config struct {
	Namespace            string
	MetricsAddr          string
	HealthProbeAddr      string
	EnableLeaderElection bool
	ZapOptions           zap.Options
}
