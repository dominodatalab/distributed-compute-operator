package logging

import (
	"context"

	"github.com/go-logr/logr"
)

// ContextLogger extends logr.Logger with the ability to store and retrieve
// custom loggers using contexts. This interface is intended for use with
// request-scoped operations.
type ContextLogger interface {
	logr.Logger

	NewContext(ctx context.Context, keysAndValues ...interface{}) (context.Context, logr.Logger)
	FromContext(context.Context) logr.Logger
}

type loggerKeyType int

const loggerKey loggerKeyType = iota

type logger struct {
	logr.Logger
}

// New returns a ContextLogger with an embedded logr.Logger instance. It can be
// used just like a regular logger when its extra context-based capabilities
// are not required.
func New(log logr.Logger) ContextLogger {
	return &logger{log}
}

// NewContext returns a copy of the parent context with a new logger that has
// the key-value pairs added to it.
func (l *logger) NewContext(ctx context.Context, keysAndValues ...interface{}) (context.Context, logr.Logger) {
	log := l.WithValues(keysAndValues...)
	return context.WithValue(ctx, loggerKey, log), log
}

// FromContext returns a logger from the context. The root logger is returned
// if the context is nil or does not contain a logger.
func (l *logger) FromContext(ctx context.Context) logr.Logger {
	if ctx == nil {
		return l
	}
	if ctxLog, ok := ctx.Value(loggerKey).(logr.Logger); ok {
		return ctxLog
	}

	return l
}
