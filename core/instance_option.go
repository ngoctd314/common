package core

import "time"

func WithLogger(logger Logger) InstanceOption {
	return func(i *Instance) {
		if logger != nil {
			i.logger = logger
		}
	}
}

func WithGracefulShutdown(timeout time.Duration) InstanceOption {
	return func(i *Instance) {
		if timeout > 0 {
			i.shutdownTimeout = timeout
		}
	}
}
