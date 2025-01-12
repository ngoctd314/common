package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"syscall"
	"testing"
	"time"
)

type mockApp struct {
	startFunc    func()
	shutdownFunc func() error
	shutdownErr  error
}

func (a *mockApp) Start(ctx context.Context) {
	if a.startFunc != nil {
		a.startFunc()
	}
}
func (a *mockApp) Shutdown(ctx context.Context) error {
	if a.shutdownFunc != nil {
		return a.shutdownFunc()
	}

	return a.shutdownErr
}

func Test_Instance_BootStrap(t *testing.T) {
	t.Parallel()

	type args struct {
		app  App
		opts []InstanceOption
	}

	bufLogger := bytes.NewBuffer(nil)

	opts := []InstanceOption{
		WithGracefulShutdown(time.Second),
		WithLogger(slog.New(slog.NewTextHandler(bufLogger, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == "time" {
					return slog.Attr{}
				}
				return slog.Attr{Key: a.Key, Value: a.Value}
			},
		}))),
	}

	testCases := []struct {
		name    string
		args    args
		wantLog []string
	}{
		{
			name: "test graceful shutdown with syscall SIGINT",
			args: args{
				app: &mockApp{
					startFunc: func() {
						syscall.Kill(syscall.Getpid(), syscall.SIGINT)
					},
				},
				opts: opts,
			},
			wantLog: []string{
				"level=WARN msg=\"receive os.Signal: interrupt\"",
				"level=INFO msg=\"shutdown complete after",
			},
		},
		{
			name: "test graceful shutdown with syscall SIGTERM",
			args: args{
				app: &mockApp{
					startFunc: func() {
						syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
					},
				},
				opts: opts,
			},
			wantLog: []string{
				"level=WARN msg=\"receive os.Signal: terminated\"",
				"level=INFO msg=\"shutdown complete after",
			},
		},
		{
			name: "test graceful shutdown when app panic",
			args: args{
				app: &mockApp{
					startFunc: func() {
						panic("st went wrong")
					},
				},
				opts: opts,
			},
			wantLog: []string{
				"level=ERROR msg=recover reason=\"st went wrong\"",
				"level=INFO msg=\"shutdown complete after",
			},
		},
		{
			name: "test error when shutdown",
			args: args{
				app: &mockApp{
					startFunc: func() {
						syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
					},
					shutdownErr: fmt.Errorf("st went wrong"),
				},
				opts: opts,
			},
			wantLog: []string{
				"level=WARN msg=\"receive os.Signal: terminated\"",
				"level=ERROR msg=\"error occur when Shutdown\"",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			instance := NewInstance(context.Background(), tc.args.app, tc.args.opts...)
			instance.Bootstrap()

			for _, v := range tc.wantLog {
				out, err := bufLogger.ReadString('\n')
				if err != nil && err != io.EOF {
					t.Errorf("failed to read log, err: %v", err)
				}
				if !strings.Contains(out, v) {
					t.Errorf("want log: %s, got: %s", v, out)
				}
			}
		})
	}
}
