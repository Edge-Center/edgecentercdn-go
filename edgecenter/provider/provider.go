package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	c := &Client{httpc: httpc, baseURL: strings.TrimSuffix(baseURL, "/")}

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

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+"/"+strings.TrimPrefix(path, "/"), body)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		apiErr := edgecenter.NewAPIError(resp.StatusCode, sentinelForStatusCode(resp.StatusCode))

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read error response %d: %w", resp.StatusCode, err)
		}

		if len(body) > 0 {
			if err := json.Unmarshal(body, apiErr); err != nil {
				apiErr.Message = strings.TrimSpace(string(body))
			}
		}

		return fmt.Errorf("%s %s: %w", method, path, apiErr)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decode successful resp %d: %w", resp.StatusCode, err)
		}
	}

	return nil
}

func sentinelForStatusCode(statusCode int) error {
	switch statusCode {
	case http.StatusBadRequest:
		return edgecenter.ErrBadRequest
	case http.StatusUnauthorized:
		return edgecenter.ErrUnauthorized
	case http.StatusForbidden:
		return edgecenter.ErrForbidden
	case http.StatusNotFound:
		return edgecenter.ErrNotFound
	case http.StatusConflict:
		return edgecenter.ErrConflict
	case http.StatusTooManyRequests:
		return edgecenter.ErrRateLimit
	default:
		return nil
	}
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
