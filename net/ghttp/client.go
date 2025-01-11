package ghttp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/ngoctd314/common/env"
)

type Client interface {
	Do(r *http.Request) (*http.Response, error)
	Client() *http.Client
}

type httpClient struct {
	client *http.Client
}

func NewClient(opts ...clientOption) *httpClient {
	instance := &http.Client{
		Timeout: env.GetWithDefault("http.client.timeout", time.Second*10),
	}

	client := &httpClient{
		client: instance,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *httpClient) Do(r *http.Request) (*http.Response, error) {
	return c.client.Do(r)
}

func (c *httpClient) Client() *http.Client {
	return c.client
}

func JSONReader(data any) io.Reader {
	b, _ := json.Marshal(data)
	return bytes.NewReader(b)
}
