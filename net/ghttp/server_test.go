package ghttp

import (
	"errors"
	"net/http"
	"os"
	"testing"
)

func TestNewServer(t *testing.T) {
	t.Parallel()

	type args struct {
		handler http.Handler
	}

	type hook struct {
		before func()
		after  func()
	}

	after := func() {
		os.Setenv("HTTP_SERVER_ADDR", "")
		os.Setenv("HTTP_SERVER_CFG", "")
	}

	testCases := []struct {
		name    string
		args    args
		hook    hook
		wantErr error
	}{
		{
			name: "test NewServer success",
			args: args{
				handler: http.DefaultServeMux,
			},
			hook: hook{
				before: func() {
					os.Setenv("HTTP_SERVER_ADDR", ":8080")
					os.Setenv("HTTP_SERVER_CFG", "readHeaderTimeout=10s&readTimeout=10s&writeTimeout=10s&idleTimeout=10s&maxHeaderBytes=1000")
				},
				after: after,
			},
		},
		{
			name: "test NewServer invalid handler",
			args: args{
				handler: nil,
			},
			wantErr: errInvalidHandler,
		},
		{
			name: "test NewServer invalid readHeaderTimeout",
			args: args{
				handler: http.DefaultServeMux,
			},
			hook: hook{
				before: func() {
					os.Setenv("HTTP_SERVER_ADDR", ":8080")
					os.Setenv("HTTP_SERVER_CFG", "readTimeout=10s&writeTimeout=10s&idleTimeout=10s&maxHeaderBytes=1000")
				},
				after: after,
			},
			wantErr: errInvalidServerCfg,
		},
		{
			name: "test NewServer invalid readTimeout",
			args: args{
				handler: http.DefaultServeMux,
			},
			hook: hook{
				before: func() {
					os.Setenv("HTTP_SERVER_ADDR", ":8080")
					os.Setenv("HTTP_SERVER_CFG", "readHeaderTimeout=10s&writeTimeout=10s&idleTimeout=10s&maxHeaderBytes=1000")
				},
				after: after,
			},
			wantErr: errInvalidServerCfg,
		},
		{
			name: "test NewServer invalid writeTimeout",
			args: args{
				handler: http.DefaultServeMux,
			},
			hook: hook{
				before: func() {
					os.Setenv("HTTP_SERVER_ADDR", ":8080")
					os.Setenv("HTTP_SERVER_CFG", "readHeaderTimeout=10s&readTimeout=10s&idleTimeout=10s&maxHeaderBytes=1000")
				},
				after: after,
			},
			wantErr: errInvalidServerCfg,
		},
		{
			name: "test NewServer invalid idleTimeout",
			args: args{
				handler: http.DefaultServeMux,
			},
			hook: hook{
				before: func() {
					os.Setenv("HTTP_SERVER_ADDR", ":8080")
					os.Setenv("HTTP_SERVER_CFG", "readHeaderTimeout=10s&readTimeout=10s&writeTimeout=10s&maxHeaderBytes=1000")
				},
				after: after,
			},
			wantErr: errInvalidServerCfg,
		},
		{
			name: "test NewServer invalid maxHeaderBytes",
			args: args{
				handler: http.DefaultServeMux,
			},
			hook: hook{
				before: func() {
					os.Setenv("HTTP_SERVER_ADDR", ":8080")
					os.Setenv("HTTP_SERVER_CFG", "readHeaderTimeout=10s&readTimeout=10s&writeTimeout=10s&idleTimeout=10s")
				},
				after: after,
			},
			wantErr: errInvalidServerCfg,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.hook.before != nil {
				tc.hook.before()
			}
			if tc.hook.after != nil {
				defer tc.hook.after()
			}

			_, err := NewServer(tc.args.handler)
			if (err != nil || tc.wantErr != nil) && !errors.Is(err, tc.wantErr) {
				t.Fatalf("NewServer() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
