package core

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Context struct {
	context.Context

	Log      logr.Logger
	Object   client.Object
	Client   client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Patch    *Patch
}
