package ghttp

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
)

type serverOption func(s *Server)

// WithServerLogger allow use your own logger style
// if it is not set or equal to nil, slog.Default is used
func WithServerLogger(l Logger) serverOption {
	return func(s *Server) {
		if l != nil {
			s.logger = l
		}
	}
}

// DisableGeneralOptionsHandler, if true, passes "OPTIONS *" requests to the Handler,
// otherwise responds with 200 OK and Content-Length: 0.
func WithDisableGeneralOptionsHandler(diable bool) serverOption {
	return func(s *Server) {
		s.instance.DisableGeneralOptionsHandler = diable
	}
}

// MaxHeaderBytes controls the maximum number of bytes the
// server will read parsing the request header's keys and
// values, including the request line. It does not limit the
// size of the request body.
// If zero, DefaultMaxHeaderBytes is used.
func WithMaxHeaderBytes(maxBytes int) serverOption {
	return func(s *Server) {
		if maxBytes > 0 {
			s.instance.MaxHeaderBytes = maxBytes
		}
	}
}

// TLSConfig optionally provides a TLS configuration for use
// by ServeTLS and ListenAndServeTLS. Note that this value is
// cloned by ServeTLS and ListenAndServeTLS, so it's not
// possible to modify the configuration with methods like
// tls.Config.SetSessionTicketKeys. To use
// SetSessionTicketKeys, use Server.Serve with a TLS Listener
// instead.
func WithTLSConfig(cnf *tls.Config) serverOption {
	return func(s *Server) {
		if cnf != nil {
			s.instance.TLSConfig = cnf
		}
	}
}

// TLSNextProto optionally specifies a function to take over
// ownership of the provided TLS connection when an ALPN
// protocol upgrade has occurred. The map key is the protocol
// name negotiated. The Handler argument should be used to
// handle HTTP requests and will initialize the Request's TLS
// and RemoteAddr if not already set. The connection is
// automatically closed when the function returns.
// If TLSNextProto is not nil, HTTP/2 support is not enabled
// automatically.
func WithTLSNextProto(tlsNextProto map[string]func(*http.Server, *tls.Conn, http.Handler)) serverOption {
	return func(s *Server) {
		s.instance.TLSNextProto = tlsNextProto
	}
}

// ConnState specifies an optional callback function that is
// called when a client connection changes state. See the
// ConnState type and associated constants for details.
func WithConnState(connState func(net.Conn, http.ConnState)) serverOption {
	return func(s *Server) {
		s.instance.ConnState = connState
	}
}

// BaseContext optionally specifies a function that returns
// the base context for incoming requests on this server.
// The provided Listener is the specific Listener that's
// about to start accepting requests.
// If BaseContext is nil, the default is context.Background().
// If non-nil, it must return a non-nil context.
func WithBaseContext(baseContext func(net.Listener) context.Context) serverOption {
	return func(s *Server) {
		s.instance.BaseContext = baseContext
	}
}

// ConnContext optionally specifies a function that modifies
// the context used for a new connection c. The provided ctx
// is derived from the base context and has a ServerContextKey
// value.
func WithConnContext(connContext func(ctx context.Context, c net.Conn) context.Context) serverOption {
	return func(s *Server) {
		s.instance.ConnContext = connContext
	}
}

// WithHandlerTimeout use http.TimeoutHandler to handle request timeout
// if handlerTimeout is zero, it will be ignored
// WARN: sse, websocket, long-polling, etc... should not use this option
// func WithHandlerTimeout(handlerTimeout time.Duration) serverOption {
// 	return func(s *Server) {
// 		if handlerTimeout > 0 {
// 			s.instance.Handler = http.TimeoutHandler(s.instance.Handler, handlerTimeout, "server timeout")
// 		}
// 	}
// }
