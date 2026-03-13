package sslcerts

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestSSLCertService_Create(t *testing.T) {
	expected := Cert{
		ID:            1,
		Name:          "test-cert",
		CertIssuer:    "Let's Encrypt",
		CertSubjectCN: "example.com",
	}

	ts := setupMockServer(t, http.MethodPost, "/cdn/ssl/certificates", http.StatusOK, expected)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.Create(context.Background(), &CreateRequest{
		Name:       "test-cert",
		Cert:       "-----BEGIN CERTIFICATE-----\nMIIB...",
		PrivateKey: "-----BEGIN PRIVATE KEY-----\nMIIE...",
	})

	require.NoError(t, err)
	assert.Equal(t, &expected, result)
}

func TestSSLCertService_Get(t *testing.T) {
	notBefore := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	notAfter := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	expected := Cert{
		ID:                1,
		Name:              "test-cert",
		CertIssuer:        "Let's Encrypt",
		CertSubjectCN:     "example.com",
		ValidityNotBefore: notBefore,
		ValidityNotAfter:  notAfter,
	}

	ts := setupMockServer(t, http.MethodGet, "/cdn/ssl/certificates/1", http.StatusOK, expected)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.Get(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "test-cert", result.Name)
	assert.Equal(t, "Let's Encrypt", result.CertIssuer)
}

func TestSSLCertService_Delete(t *testing.T) {
	ts := setupMockServer(t, http.MethodDelete, "/cdn/ssl/certificates/1", http.StatusOK, nil)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	err := service.Delete(context.Background(), 1)

	require.NoError(t, err)
}

func TestSSLCertService_NotFound(t *testing.T) {
	ts := setupMockServer(t, http.MethodGet, "/cdn/ssl/certificates/999", http.StatusNotFound, nil)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	_, err := service.Get(context.Background(), 999)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestSSLCertService_BadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"Message": "invalid certificate"})
	}))
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	_, err := service.Create(context.Background(), &CreateRequest{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid certificate")
}
