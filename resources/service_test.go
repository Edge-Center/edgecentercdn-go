package resources

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

func TestResourceService_Create(t *testing.T) {
	expected := Resource{
		ID:     1,
		Cname:  "cdn.example.com",
		Status: ActiveResourceStatus,
		Active: true,
		Client: 12345,
	}

	ts := setupMockServer(t, http.MethodPost, "/cdn/resources", http.StatusOK, expected)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.Create(context.Background(), &CreateRequest{
		Cname:          "cdn.example.com",
		OriginProtocol: HTTPProtocol,
		Origin:         "origin.example.com",
	})

	require.NoError(t, err)
	assert.Equal(t, &expected, result)
}

func TestResourceService_Get(t *testing.T) {
	expected := Resource{
		ID:     1,
		Cname:  "cdn.example.com",
		Status: ActiveResourceStatus,
		Active: true,
		Client: 12345,
	}

	ts := setupMockServer(t, http.MethodGet, "/cdn/resources/1", http.StatusOK, expected)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.Get(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, &expected, result)
}

func TestResourceService_Update(t *testing.T) {
	expected := Resource{
		ID:     1,
		Cname:  "cdn.example.com",
		Status: ActiveResourceStatus,
		Active: true,
	}

	ts := setupMockServer(t, http.MethodPut, "/cdn/resources/1", http.StatusOK, expected)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.Update(context.Background(), 1, &UpdateRequest{Active: true})

	require.NoError(t, err)
	assert.Equal(t, &expected, result)
}

func TestResourceService_Delete(t *testing.T) {
	ts := setupMockServer(t, http.MethodDelete, "/cdn/resources/1", http.StatusOK, nil)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	err := service.Delete(context.Background(), 1)

	require.NoError(t, err)
}

func TestResourceService_Page(t *testing.T) {
	expected := PaginatedResource{
		Count: 12,
		Results: []Resource{
			{
				ID:          1,
				Status:      ActiveResourceStatus,
				Active:      true,
				Client:      12345,
				OriginGroup: 1,
				Cname:       "cdn1.example.com",
			},
			{
				ID:          2,
				Status:      SuspendedResourceStatus,
				Active:      false,
				Client:      12345,
				OriginGroup: 2,
				Cname:       "cdn2.example.com",
			},
		},
	}

	ts := setupMockServer(t, http.MethodGet, "/cdn/resources", http.StatusOK, expected)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.Page(context.Background(), 10, 10, &ListFilterRequest{
		Fields: []string{"id", "cname"},
		Status: []ResourceStatus{ActiveResourceStatus, ProcessedResourceStatus},
	})

	require.NoError(t, err)
	assert.Equal(t, &expected, result)
}

func TestResourceService_List(t *testing.T) {
	expected := []Resource{
		{
			ID:          1,
			Status:      ActiveResourceStatus,
			Active:      true,
			Client:      12345,
			OriginGroup: 1,
			Cname:       "cdn1.example.com",
		},
		{
			ID:          2,
			Status:      SuspendedResourceStatus,
			Active:      false,
			Client:      12345,
			OriginGroup: 2,
			Cname:       "cdn2.example.com",
		},
	}

	ts := setupMockServer(t, http.MethodGet, "/cdn/resources", http.StatusOK, expected)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.List(context.Background(), &ListFilterRequest{
		Fields: []string{"id", "cname"},
		Status: []ResourceStatus{ActiveResourceStatus, ProcessedResourceStatus},
	})

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestResourceService_Count(t *testing.T) {
	ts := setupMockServer(t, http.MethodGet, "/cdn/resources", http.StatusOK, PaginatedResource{Count: 42})
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	count, err := service.Count(context.Background(), nil)

	require.NoError(t, err)
	assert.Equal(t, uint(42), count)
}

func TestResourceService_Create_Error(t *testing.T) {
	ts := setupMockServer(t, http.MethodPost, "/cdn/resources", http.StatusNotFound, nil)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	_, err := service.Create(context.Background(), &CreateRequest{Cname: "cdn.example.com"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "request:")
}

func TestResourceService_NotFound(t *testing.T) {
	ts := setupMockServer(t, http.MethodGet, "/cdn/resources/999", http.StatusNotFound, nil)
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	_, err := service.Get(context.Background(), 999)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestResourceService_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"Message": "server error"})
	}))
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	_, err := service.Get(context.Background(), 1)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "server error")
}

func TestResourceService_InvalidJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	_, err := service.Get(context.Background(), 1)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode")
}
