package httpx

import (
	"net"
	"net/http"
	"time"
)

type Client struct {
	*http.Client
}

func NewClient() *Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &Client{
		Client: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		},
	}
}