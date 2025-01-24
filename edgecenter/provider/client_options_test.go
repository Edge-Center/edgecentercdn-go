package provider

import (
	"context"
	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestAuthenticatedHeaders(t *testing.T) {
	type args struct {
		apiKey string
	}
	tests := []struct {
		name  string
		args  args
		wantM map[string]string
	}{
		{
			name:  "valid API key",
			args:  args{apiKey: "test-api-key"},
			wantM: map[string]string{"Authorization": "APIKey test-api-key"},
		},
		{
			name:  "empty API key",
			args:  args{apiKey: ""},
			wantM: map[string]string{"Authorization": "APIKey "}, // Assuming empty key is allowed
		},
		{
			name:  "whitespace API key",
			args:  args{apiKey: " "},
			wantM: map[string]string{"Authorization": "APIKey  "}, // Handling case with space key
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotM := AuthenticatedHeaders(tt.args.apiKey); !reflect.DeepEqual(gotM, tt.wantM) {
				t.Errorf("AuthenticatedHeaders() = %v, want %v", gotM, tt.wantM)
			}
		})
	}
}

func TestWithSigner(t *testing.T) {
	tests := []struct {
		name         string
		signer       edgecenter.RequestSignerFunc
		expectedAuth string
		expectedErr  bool
		method       string
		path         string
		payload      interface{}
	}{
		{
			name:         "Valid signer with correct Authorization",
			signer:       edgecenter.RequestSignerFunc(func(req *http.Request) error { req.Header.Set("Authorization", "APIKey test-api-key"); return nil }),
			expectedAuth: "APIKey test-api-key",
			expectedErr:  false,
			method:       http.MethodGet,
			path:         "/test",
			payload:      nil,
		},
		{
			name:         "Valid signer with different Authorization",
			signer:       edgecenter.RequestSignerFunc(func(req *http.Request) error { req.Header.Set("Authorization", "APIKey different-api-key"); return nil }),
			expectedAuth: "APIKey different-api-key",
			expectedErr:  false,
			method:       http.MethodGet,
			path:         "/test2",
			payload:      nil,
		},
		{
			name:         "Invalid signer with Unauthorized response",
			signer:       edgecenter.RequestSignerFunc(func(req *http.Request) error { req.Header.Set("Authorization", "APIKey no-api-key"); return nil }),
			expectedAuth: "APIKey no-api-key",
			expectedErr:  true,
			method:       http.MethodGet,
			path:         "/test3",
			payload:      nil,
		},
		{
			name:         "Valid signer with POST method and payload",
			signer:       edgecenter.RequestSignerFunc(func(req *http.Request) error { req.Header.Set("Authorization", "APIKey valid-api-key"); return nil }),
			expectedAuth: "APIKey valid-api-key",
			expectedErr:  false,
			method:       http.MethodPost,
			path:         "/test4",
			payload:      map[string]string{"key": "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.expectedAuth, r.Header.Get("Authorization"))
				if tt.expectedErr {
					w.WriteHeader(http.StatusUnauthorized)
				} else {
					w.WriteHeader(http.StatusOK)
				}
			}))
			defer ts.Close()

			client := NewClient(ts.URL, WithSigner(tt.signer))

			var err error
			if tt.payload != nil {
				err = client.Request(context.Background(), tt.method, tt.path, tt.payload, nil)
			} else {
				err = client.Request(context.Background(), tt.method, tt.path, nil, nil)
			}

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	tests := []struct {
		name           string
		timeout        time.Duration
		serverResponse time.Duration
		expectErr      bool
	}{
		{
			name:           "Valid timeout - within server response time",
			timeout:        2 * time.Second,
			serverResponse: 1 * time.Second,
			expectErr:      false,
		},
		{
			name:           "Valid timeout - server response exceeds timeout",
			timeout:        2 * time.Second,
			serverResponse: 3 * time.Second,
			expectErr:      true,
		},
		{
			name:           "Large timeout",
			timeout:        3 * time.Second,
			serverResponse: 1 * time.Second,
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(tt.serverResponse)
				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()

			client := NewClient(ts.URL, WithTimeout(tt.timeout))

			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout+1*time.Second)
			defer cancel()

			err := client.Request(ctx, "GET", "/", nil, nil)

			if tt.expectErr {
				assert.Error(t, err, "expected an error but got nil")
			} else {
				assert.NoError(t, err, "unexpected error")
			}

			if client.httpc.Timeout != tt.timeout {
				t.Errorf("expected timeout %v, got %v", tt.timeout, client.httpc.Timeout)
			}
		})
	}
}

func TestWithUA(t *testing.T) {
	tests := []struct {
		name       string
		ua         string
		wantHeader string
	}{
		{
			name:       "Valid User-Agent",
			ua:         "MyApp/1.0",
			wantHeader: "MyApp/1.0",
		},
		{
			name:       "Another User-Agent",
			ua:         "CustomUA/2.0",
			wantHeader: "CustomUA/2.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.Header.Get("User-Agent"); got != tt.wantHeader {
					t.Errorf("User-Agent = %v, want %v", got, tt.wantHeader)
				}
			}))
			defer ts.Close()

			client := NewClient(ts.URL, WithUA(tt.ua))

			err := client.Request(context.Background(), http.MethodGet, "/", nil, nil)
			if err != nil {
				t.Errorf("Request failed: %v", err)
			}
		})
	}
}
