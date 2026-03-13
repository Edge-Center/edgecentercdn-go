package lecerts

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Edge-Center/edgecentercdn-go/edgecenter/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockServer(t *testing.T, method, path string, statusCode int, response interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, method, r.Method)
		assert.Equal(t, path, r.URL.Path)
		w.WriteHeader(statusCode)
		if response != nil {
			json.NewEncoder(w).Encode(response)
		}
	}))
}

func TestLECertService_GetLECert(t *testing.T) {
	expected := LECertStatus{
		ID:       1,
		Active:   true,
		Resource: 100,
		Started:  "2024-01-01T00:00:00Z",
		Statuses: []LEStatusDetail{
			{Status: "done", Created: "2024-01-01T00:00:00Z"},
		},
	}

	ts := setupMockServer(t, http.MethodGet, "/cdn/resources/100/ssl/le/status", http.StatusOK, expected)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.GetLECert(context.Background(), 100)

	require.NoError(t, err)
	assert.Equal(t, &expected, result)
}

func TestLECertService_CreateLECert(t *testing.T) {
	ts := setupMockServer(t, http.MethodPost, "/cdn/resources/100/ssl/le/issue", http.StatusOK, nil)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	err := service.CreateLECert(context.Background(), 100)

	require.NoError(t, err)
}

func TestLECertService_UpdateLECert(t *testing.T) {
	ts := setupMockServer(t, http.MethodPost, "/cdn/resources/100/ssl/le/renew", http.StatusOK, nil)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	err := service.UpdateLECert(context.Background(), 100)

	require.NoError(t, err)
}

func TestLECertService_DeleteLECert(t *testing.T) {
	tests := []struct {
		name  string
		force bool
	}{
		{name: "force=true", force: true},
		{name: "force=false", force: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/cdn/resources/100/ssl/le/revoke", r.URL.Path)
				assert.Equal(t, tt.force, r.URL.Query().Get("force") == "true")
				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()

			service := NewService(provider.NewClient(ts.URL))
			err := service.DeleteLECert(context.Background(), 100, tt.force)

			require.NoError(t, err)
		})
	}
}

func TestLECertService_CancelLECert(t *testing.T) {
	tests := []struct {
		name   string
		active bool
	}{
		{name: "active=true", active: true},
		{name: "active=false", active: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				assert.Equal(t, "/cdn/resources/100/ssl/le/status", r.URL.Path)

				body, err := io.ReadAll(r.Body)
				if !assert.NoError(t, err) {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				var payload map[string]bool
				if !assert.NoError(t, json.Unmarshal(body, &payload)) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				assert.Equal(t, tt.active, payload["active"])

				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()

			service := NewService(provider.NewClient(ts.URL))
			err := service.CancelLECert(context.Background(), 100, tt.active)

			require.NoError(t, err)
		})
	}
}

func TestLECertService_NotFound(t *testing.T) {
	ts := setupMockServer(t, http.MethodGet, "/cdn/resources/999/ssl/le/status", http.StatusNotFound, nil)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	_, err := service.GetLECert(context.Background(), 999)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestLECertService_Conflict(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"Message": "certificate already issued"})
	}))
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	err := service.CreateLECert(context.Background(), 100)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "certificate already issued")
}
