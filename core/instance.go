package core

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Instance struct {
	baseCtx         context.Context
	cancelFunc      func()
	app             App
	logger          Logger
	shutdownTimeout time.Duration // graceful shutdown timeout
}

type App interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context) error
}
type InstanceOption func(*Instance)

func NewInstance(ctx context.Context, app App, opts ...InstanceOption) *Instance {
	baseCtx, cancel := context.WithCancel(ctx)

	instance := &Instance{
		baseCtx:         baseCtx,
		cancelFunc:      cancel,
		app:             app,
		logger:          slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		shutdownTimeout: time.Second * 10, // the default graceful shutdown timeout is 10s
	}

	for _, opt := range opts {
		opt(instance)
	}

	return instance
}

// Bootstrap start app and handle graceful shutdown
func (i *Instance) Bootstrap() {
	// start instance
	go func() {
		defer func() {
			if r := recover(); r != nil {
				i.logger.Error("recover", "reason", r)
				i.cancelFunc()
			}
		}()
		i.app.Start(i.baseCtx)
	}()

	// handle graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// wait for termination signal or context done
	select {
	case v := <-signalCh:
		i.logger.Warn(fmt.Sprintf("receive os.Signal: %s", v))
	case <-i.baseCtx.Done():
	}

	now := time.Now()

	shutdownCtx, cancel := context.WithTimeout(i.baseCtx, i.shutdownTimeout)
	defer cancel()

	// create a WaitGroup to keep track of shutdown goroutines
	if err := i.app.Shutdown(shutdownCtx); err != nil {
		i.logger.Error("error occur when Shutdown", "err", err)
		return
	}
	i.logger.Info(fmt.Sprintf("shutdown complete after %f seconds", time.Since(now).Seconds()))
}
