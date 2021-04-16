package http

import (
	h "net/http"
)

type HttpClient struct {
	Client  h.Client
	Request h.Request
}

func New() *HttpClient {
	client := h.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}

	return &HttpClient{
		Client: client,
	}
}

func (c *HttpClient) SendRequest(req h.Request) (*h.Response, error) {
	return c.Client.Do(&req)
}
