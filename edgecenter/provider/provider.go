package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
)

type Client struct {
	httpc   *http.Client
	signer  edgecenter.RequestSigner
	ua      string
	baseURL string
}

func NewClient(baseURL string, opts ...ClientOption) *Client {
	httpc := &http.Client{Timeout: time.Minute}
	c := &Client{httpc: httpc, baseURL: baseURL}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) Request(ctx context.Context, method, path string, payload interface{}, result interface{}) error {
	var body io.Reader
	if payload != nil {
		payloadBuf := new(bytes.Buffer)
		if err := json.NewEncoder(payloadBuf).Encode(payload); err != nil {
			return fmt.Errorf("encode req payload: %w", err)
		}

		body = payloadBuf
	}

	// TODO: figure out how to drop trailing slash
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return fmt.Errorf("resource not found at path: %s", path)
		case http.StatusUnauthorized:
			return fmt.Errorf("unauthorized. Invalid or expired credentials, please check your authentication setup")
		default:
			var errResp edgecenter.ErrorResponse
			if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
				return fmt.Errorf("decode err resp %d: %w", resp.StatusCode, err)
			}
			return &errResp
		}
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decode successful resp %d: %w", resp.StatusCode, err)
		}
	}

	return nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")

	if c.ua != "" {
		req.Header.Set("User-Agent", c.ua)
	}

	if c.signer != nil {
		if err := c.signer.Sign(req); err != nil {
			return nil, err
		}
	}

	return c.httpc.Do(req)
}
