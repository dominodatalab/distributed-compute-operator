package manager

import (
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	istioscheme "istio.io/client-go/pkg/clientset/versioned/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/controllers"
	//+kubebuilder:scaffold:imports
)

const leaderElectionID = "a846cbf2.dominodatalab.com"

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

// Start creates a new controller manager, configures and registers all
// reconcilers/webhooks with the manager, and starts their control loops.
func Start(cfg *controllers.Config) error {
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&cfg.ZapOptions)))

	mgrOpts := manager.Options{
		Scheme: scheme,
		Metrics: server.Options{
			BindAddress: cfg.MetricsAddr,
		},
		WebhookServer: webhook.NewServer(webhook.Options{
			Port: cfg.WebhookServerPort,
		}),
		HealthProbeBindAddress: cfg.HealthProbeAddr,
		LeaderElection:         cfg.EnableLeaderElection,
		LeaderElectionID:       leaderElectionID,
	}

	if len(cfg.Namespaces) > 0 {
		setupLog.Info("Limiting reconciliation watch", "namespaces", cfg.Namespaces)
		// TODO: double check this - MultiNamespacedCacheBuilder is removed
		// mgrOpts.NewCache = cache.MultiNamespacedCacheBuilder(cfg.Namespaces)
		// https://github.com/kubernetes-sigs/controller-runtime/pull/2157
		namespaces := map[string]cache.Config{}
		for _, n := range cfg.Namespaces {
			namespaces[n] = cache.Config{}
		}
		mgrOpts.Cache = cache.Options{
			DefaultNamespaces: namespaces,
		}
	} else {
		setupLog.Info("Watching all namespaces")
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), mgrOpts)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		return err
	}

	if err = mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		return err
	}
	if err = mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		return err
	}

	enableWebHooks := os.Getenv("ENABLE_WEBHOOKS") != "false" // TODO: add to config

	for _, builder := range controllers.BuilderFuncs {
		if err = builder(mgr, enableWebHooks, cfg); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", builder)
			return err
		}
	}

	// NOTE: old approach to setup

	if err = (&controllers.RayClusterReconciler{
		Client:       mgr.GetClient(),
		Log:          ctrl.Log.WithName("controllers").WithName("RayCluster"),
		Scheme:       mgr.GetScheme(),
		IstioEnabled: cfg.IstioEnabled,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "RayCluster")
		return err
	}

	if err = (&controllers.SparkClusterReconciler{
		Client:       mgr.GetClient(),
		Log:          ctrl.Log.WithName("controllers").WithName("SparkCluster"),
		Scheme:       mgr.GetScheme(),
		IstioEnabled: cfg.IstioEnabled,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "SparkCluster")
		return err
	}

	if enableWebHooks {
		if err = (&dcv1alpha1.RayCluster{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "RayCluster")
			return err
		}
		if err = (&dcv1alpha1.SparkCluster{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "SparkCluster")
			return err
		}
	}

	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		return err
	}

	return nil
}

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(dcv1alpha1.AddToScheme(scheme))
	utilruntime.Must(istioscheme.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}
