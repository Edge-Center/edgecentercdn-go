package provider

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Request(t *testing.T) {
	type fields struct {
		httpc   *http.Client
		signer  edgecenter.RequestSigner
		ua      string
		baseURL string
	}
	type args struct {
		ctx     context.Context
		method  string
		path    string
		payload interface{}
		result  interface{}
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		baseURLTrailing bool
		wantErr         bool
	}{
		{
			name: "successful request with valid response",
			fields: fields{
				httpc: &http.Client{},
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				path:   "/test",
				result: &map[string]interface{}{"key": "value"},
			},
			wantErr: false,
		},
		{
			name: "request with server error",
			fields: fields{
				httpc: &http.Client{},
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				path:   "/test",
				result: &map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name: "request with custom user agent",
			fields: fields{
				httpc: &http.Client{},
				ua:    "custom-user-agent",
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				path:   "/test",
				result: &map[string]interface{}{"key": "value"},
			},
			wantErr: false,
		},
		{
			name: "request with timeout error",
			fields: fields{
				httpc: &http.Client{
					Timeout: 1 * time.Millisecond,
				},
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				path:   "/test",
				result: &map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name: "request with valid payload",
			fields: fields{
				httpc: &http.Client{},
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodPost,
				path:   "/test",
				payload: map[string]string{
					"key": "value",
				},
				result: &map[string]interface{}{"response": "success"},
			},
			wantErr: false,
		},
		{
			name: "request with invalid path",
			fields: fields{
				httpc: &http.Client{},
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				path:   "/invalid-path",
				result: &map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name: "baseURL with trailing slash and path with leading slash produce single slash",
			fields: fields{
				httpc: &http.Client{},
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				path:   "/test",
				result: &map[string]interface{}{"key": "value"},
			},
			baseURLTrailing: true,
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.Contains(r.RequestURI, "//") {
					t.Errorf("double slash in request URI: %s", r.RequestURI)
				}

				if r.URL.Path != tt.args.path {
					t.Errorf("Expected path %s, but got %s", tt.args.path, r.URL.Path)
				}

				if tt.args.payload != nil {
					decoder := json.NewDecoder(r.Body)
					var payload map[string]string
					if err := decoder.Decode(&payload); err != nil {
						t.Errorf("Failed to decode payload: %v", err)
					}
				}

				if tt.wantErr {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				} else {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(tt.args.result)
				}
			}))
			defer ts.Close()

			if tt.baseURLTrailing {
				tt.fields.baseURL = ts.URL + "/"
			} else {
				tt.fields.baseURL = ts.URL
			}

			var c *Client
			if tt.baseURLTrailing {
				c = NewClient(tt.fields.baseURL)
			} else {
				c = &Client{
					httpc:   tt.fields.httpc,
					signer:  tt.fields.signer,
					ua:      tt.fields.ua,
					baseURL: tt.fields.baseURL,
				}
			}

			err := c.Request(tt.args.ctx, tt.args.method, tt.args.path, tt.args.payload, tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("Request() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Request_NotFoundMappedToSentinel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"Message": "resource does not exist",
		})
	}))
	defer ts.Close()

	client := NewClient(ts.URL)

	err := client.Request(context.Background(), http.MethodGet, "/test", nil, nil)
	require.Error(t, err)
	assert.True(t, errors.Is(err, edgecenter.ErrNotFound))

	var apiErr *edgecenter.APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	assert.Equal(t, "resource does not exist", apiErr.Message)
}

func TestClient_Request_BadRequestMappedToSentinel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"Message": "invalid request",
		})
	}))
	defer ts.Close()

	client := NewClient(ts.URL)

	err := client.Request(context.Background(), http.MethodGet, "/test", nil, nil)
	require.Error(t, err)
	assert.True(t, errors.Is(err, edgecenter.ErrBadRequest))

	var apiErr *edgecenter.APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	assert.Equal(t, "invalid request", apiErr.Message)
}

func TestClient_Request_ConflictMappedToSentinel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"Message": "resource already exists",
		})
	}))
	defer ts.Close()

	client := NewClient(ts.URL)

	err := client.Request(context.Background(), http.MethodGet, "/test", nil, nil)
	require.Error(t, err)
	assert.True(t, errors.Is(err, edgecenter.ErrConflict))

	var apiErr *edgecenter.APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, http.StatusConflict, apiErr.StatusCode)
	assert.Equal(t, "resource already exists", apiErr.Message)
}

func TestClient_Request_RateLimitMappedToSentinel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"Message": "too many requests",
		})
	}))
	defer ts.Close()

	client := NewClient(ts.URL)

	err := client.Request(context.Background(), http.MethodGet, "/test", nil, nil)
	require.Error(t, err)
	assert.True(t, errors.Is(err, edgecenter.ErrRateLimit))

	var apiErr *edgecenter.APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, http.StatusTooManyRequests, apiErr.StatusCode)
	assert.Equal(t, "too many requests", apiErr.Message)
}

func TestClient_Request_PlainTextErrorResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("unauthorized"))
	}))
	defer ts.Close()

	client := NewClient(ts.URL)

	err := client.Request(context.Background(), http.MethodGet, "/test", nil, nil)
	require.Error(t, err)
	assert.True(t, errors.Is(err, edgecenter.ErrUnauthorized))

	var apiErr *edgecenter.APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
	assert.Equal(t, "unauthorized", apiErr.Message)
}

func TestClient_Request_LowercaseJSONErrorResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "access denied",
			"errors": map[string][]string{
				"token": {"invalid token"},
			},
		})
	}))
	defer ts.Close()

	client := NewClient(ts.URL)

	err := client.Request(context.Background(), http.MethodGet, "/test", nil, nil)
	require.Error(t, err)
	assert.True(t, errors.Is(err, edgecenter.ErrForbidden))

	var apiErr *edgecenter.APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, http.StatusForbidden, apiErr.StatusCode)
	assert.Equal(t, "access denied", apiErr.Message)
	assert.Len(t, apiErr.Details, 1)
	assert.Equal(t, "token", apiErr.Details[0].Field)
	assert.Equal(t, []string{"invalid token"}, apiErr.Details[0].Messages)
}
