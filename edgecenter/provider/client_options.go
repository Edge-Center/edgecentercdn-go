package provider

import (
	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
	"time"
)

type ClientOption func(*Client)

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpc.Timeout = timeout
	}
}

func WithSigner(signer edgecenter.RequestSigner) ClientOption {
	return func(c *Client) {
		c.signer = signer
	}
}

func WithSignerFunc(f edgecenter.RequestSignerFunc) ClientOption {
	return func(c *Client) {
		c.signer = f
	}
}

func WithUA(ua string) ClientOption {
	return func(c *Client) {
		c.ua = ua
	}
}
