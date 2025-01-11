package ghttp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/ngoctd314/common/env"
)

type Server struct {
	instance *http.Server
	logger   Logger
}

var (
	errInvalidHandler   = errors.New("invalid handler")
	errInvalidServerCfg = errors.New("invalid server config")
)

func NewServer(handler http.Handler, opts ...serverOption) (*Server, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if handler == nil {
		return nil, errInvalidHandler
	}

	httpServer := &http.Server{
		Handler: handler,
	}
	err := setServerCfg(httpServer, env.GetString("http.server.cfg"))
	if err != nil {
		return nil, fmt.Errorf("%w, detail: %w", errInvalidServerCfg, err)
	}

	server := &Server{
		instance: httpServer,
		logger:   logger,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server, nil
}

func (s *Server) ListenAndServe() error {
	s.logger.Info("server is listening", "addr", s.instance.Addr)
	return s.instance.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.instance.Shutdown(ctx)
}

func setServerCfg(server *http.Server, params string) error {
	q, err := url.ParseQuery(params)
	if err != nil {
		return err
	}
	server.Addr = env.GetString("http.server.addr")

	var errGroup error

	readHeaderTimeout, err := time.ParseDuration(q.Get("readHeaderTimeout"))
	server.ReadHeaderTimeout = readHeaderTimeout
	if err != nil {
		errGroup = errors.Join(errGroup, fmt.Errorf("invalid readHeaderTimeout, error: %w", err))
	}

	readTimeout, err := time.ParseDuration(q.Get("readTimeout"))
	server.ReadTimeout = readTimeout
	if err != nil {
		errGroup = errors.Join(errGroup, fmt.Errorf("invalid readTimeout, error: %w", err))
	}

	writeTimeout, err := time.ParseDuration(q.Get("writeTimeout"))
	server.WriteTimeout = writeTimeout
	if err != nil {
		errGroup = errors.Join(errGroup, fmt.Errorf("invalid writeTimeout, error: %w", err))
	}

	idleTimeout, err := time.ParseDuration(q.Get("idleTimeout"))
	server.IdleTimeout = idleTimeout
	if err != nil {
		errGroup = errors.Join(errGroup, fmt.Errorf("invalid idleTimeout, error: %w", err))
	}

	maxHeaderBytes, err := strconv.Atoi(q.Get("maxHeaderBytes"))
	if err != nil {
		errGroup = errors.Join(errGroup, fmt.Errorf("invalid maxHeaderBytes, error: %w", err))
	}
	server.MaxHeaderBytes = maxHeaderBytes

	if errGroup != nil {
		return errGroup
	}

	return nil
}
