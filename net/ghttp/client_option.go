package ghttp

import "net/http"

type clientOption func(s *httpClient)

// Transport specifies the mechanism by which individual
// HTTP requests are made.
// If nil, DefaultTransport is used.
func WithTransport(transport http.RoundTripper) clientOption {
	return func(s *httpClient) {
		if transport != nil {
			s.client.Transport = transport
		}
	}
}

// CheckRedirect specifies the policy for handling redirects.
// If CheckRedirect is not nil, the client calls it before
// following an HTTP redirect. The arguments req and via are
// the upcoming request and the requests made already, oldest
// first. If CheckRedirect returns an error, the Client's Get
// method returns both the previous Response (with its Body
// closed) and CheckRedirect's error (wrapped in a url.Error)
// instead of issuing the Request req.
// As a special case, if CheckRedirect returns ErrUseLastResponse,
// then the most recent response is returned with its body
// unclosed, along with a nil error.
//
// If CheckRedirect is nil, the Client uses its default policy,
// which is to stop after 10 consecutive requests.
func WithCheckRedirect(checkRedirect func(req *http.Request, via []*http.Request) error) clientOption {
	return func(s *httpClient) {
		if checkRedirect != nil {
			s.client.CheckRedirect = checkRedirect
		}
	}
}

// Jar specifies the cookie jar.
//
// The Jar is used to insert relevant cookies into every
// outbound Request and is updated with the cookie values
// of every inbound Response. The Jar is consulted for every
// redirect that the Client follows.
//
// If Jar is nil, cookies are only sent if they are explicitly
// set on the Request.
func WithJar(jar http.CookieJar) clientOption {
	return func(s *httpClient) {
		s.client.Jar = jar
	}
}
