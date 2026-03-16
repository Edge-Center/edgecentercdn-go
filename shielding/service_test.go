package shielding

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

func TestShieldingService_Get(t *testing.T) {
	pop := 1
	expected := ShieldingData{ShieldingPop: &pop}

	ts := setupMockServer(t, http.MethodGet, "/cdn/resources/100/shielding", http.StatusOK, expected)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.Get(context.Background(), 100)

	require.NoError(t, err)
	require.NotNil(t, result.ShieldingPop)
	assert.Equal(t, 1, *result.ShieldingPop)
}

func TestShieldingService_Update(t *testing.T) {
	pop := 2
	expected := ShieldingData{ShieldingPop: &pop}

	ts := setupMockServer(t, http.MethodPut, "/cdn/resources/100/shielding", http.StatusOK, expected)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.Update(context.Background(), 100, &UpdateShieldingData{ShieldingPop: &pop})

	require.NoError(t, err)
	require.NotNil(t, result.ShieldingPop)
	assert.Equal(t, 2, *result.ShieldingPop)
}

func TestShieldingService_GetShieldingLocations(t *testing.T) {
	expected := []ShieldingLocations{
		{ID: 1, Datacenter: "DC1", Country: "US", City: "New York"},
		{ID: 2, Datacenter: "DC2", Country: "DE", City: "Frankfurt"},
	}

	ts := setupMockServer(t, http.MethodGet, "/cdn/shielding_pop", http.StatusOK, expected)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.GetShieldingLocations(context.Background())

	require.NoError(t, err)
	require.Len(t, *result, 2)
	assert.Equal(t, "DC1", (*result)[0].Datacenter)
	assert.Equal(t, "DC2", (*result)[1].Datacenter)
}

func TestShieldingService_NotFound(t *testing.T) {
	ts := setupMockServer(t, http.MethodGet, "/cdn/resources/999/shielding", http.StatusNotFound, nil)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	_, err := service.Get(context.Background(), 999)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
