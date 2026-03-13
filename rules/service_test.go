package rules

import (
	"context"
	"encoding/json"
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

func TestRulesService_Create(t *testing.T) {
	originGroup := 10
	expected := Rule{
		ID:          1,
		Name:        "test-rule",
		Active:      true,
		Pattern:     "/images/*",
		Weight:      1,
		OriginGroup: &originGroup,
	}

	ts := setupMockServer(t, http.MethodPost, "/cdn/resources/100/locations", http.StatusOK, expected)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.Create(context.Background(), 100, &CreateRequest{
		Name:        "test-rule",
		Active:      true,
		Rule:        "/images/*",
		Weight:      1,
		OriginGroup: &originGroup,
	})

	require.NoError(t, err)
	assert.Equal(t, &expected, result)
}

func TestRulesService_Get(t *testing.T) {
	expected := Rule{
		ID:      1,
		Name:    "test-rule",
		Active:  true,
		Pattern: "/images/*",
	}

	ts := setupMockServer(t, http.MethodGet, "/cdn/resources/100/locations/1", http.StatusOK, expected)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.Get(context.Background(), 100, 1)

	require.NoError(t, err)
	assert.Equal(t, &expected, result)
}

func TestRulesService_Update(t *testing.T) {
	expected := Rule{
		ID:      1,
		Name:    "updated-rule",
		Active:  true,
		Pattern: "/videos/*",
		Weight:  2,
	}

	ts := setupMockServer(t, http.MethodPut, "/cdn/resources/100/locations/1", http.StatusOK, expected)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.Update(context.Background(), 100, 1, &UpdateRequest{
		Name:   "updated-rule",
		Active: true,
		Rule:   "/videos/*",
		Weight: 2,
	})

	require.NoError(t, err)
	assert.Equal(t, &expected, result)
}

func TestRulesService_Delete(t *testing.T) {
	ts := setupMockServer(t, http.MethodDelete, "/cdn/resources/100/locations/1", http.StatusOK, nil)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	err := service.Delete(context.Background(), 100, 1)

	require.NoError(t, err)
}

func TestRulesService_NotFound(t *testing.T) {
	ts := setupMockServer(t, http.MethodGet, "/cdn/resources/100/locations/999", http.StatusNotFound, nil)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	_, err := service.Get(context.Background(), 100, 999)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRulesService_BadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"Message": "invalid request"})
	}))
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	_, err := service.Create(context.Background(), 100, &CreateRequest{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid request")
}
