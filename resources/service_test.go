package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Edge-Center/edgecentercdn-go/edgecenter/provider"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type ListRequest struct {
	Filter *ListFilterRequest

	Offset uint
	Size   uint
}

func TestService_Page(t *testing.T) {
	tests := []struct {
		name             string
		request          *ListRequest
		response         PaginatedResource
		expectedResponse PaginatedResource
		expectedQuery    string
	}{
		{
			name: "Page resources",
			request: &ListRequest{
				Offset: 10,
				Size:   10,
				Filter: &ListFilterRequest{
					Fields: []string{"id", "cname"},
					Status: []ResourceStatus{ActiveResourceStatus, ProcessedResourceStatus},
				},
			},
			response: PaginatedResource{
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
			},
			expectedResponse: PaginatedResource{
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
			},
			expectedQuery: "fields=id%2Ccname&offset=10&size=10&status=active%2Cprocessed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/cdn/resources" || r.Method != http.MethodGet {
					http.Error(w, "not found", http.StatusNotFound)
					return
				}

				w.Header().Set("Content-Type", "application/json")

				if r.URL.RawQuery != tt.expectedQuery {
					w.WriteHeader(http.StatusBadRequest)
					errorResponse := fmt.Sprintf("invalid query parameters: expected %+v, got %+v", tt.expectedQuery, r.URL.RawQuery)
					_ = json.NewEncoder(w).Encode(map[string]string{"message": errorResponse})
					return
				}

				err := json.NewEncoder(w).Encode(tt.response)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)

					_ = json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("invalid response %s", err.Error())})
					return
				}
			})

			ts := httptest.NewServer(mockHandler)
			defer ts.Close()

			service := NewService(provider.NewClient(ts.URL))

			ctx := context.Background()

			response, err := service.Page(ctx, tt.request.Offset, tt.request.Size, tt.request.Filter)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(response, &tt.expectedResponse) {
				t.Errorf("expected %+v, got %+v", &tt.expectedResponse, response)
			}
		})
	}
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name          string
		request       *ListFilterRequest
		response      []Resource
		expected      []Resource
		expectedQuery string
	}{
		{
			name: "List resources",
			request: &ListFilterRequest{
				Fields: []string{"id", "cname"},
				Status: []ResourceStatus{ActiveResourceStatus, ProcessedResourceStatus},
			},
			response: []Resource{
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
			expected: []Resource{
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
			expectedQuery: "fields=id%2Ccname&status=active%2Cprocessed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/cdn/resources" || r.Method != http.MethodGet {
					http.Error(w, "not found", http.StatusNotFound)
					return
				}

				w.Header().Set("Content-Type", "application/json")

				if r.URL.RawQuery != tt.expectedQuery {
					w.WriteHeader(http.StatusBadRequest)
					errorResponse := fmt.Sprintf("invalid query parameters: expected %+v, got %+v", tt.expectedQuery, r.URL.RawQuery)
					_ = json.NewEncoder(w).Encode(map[string]string{"message": errorResponse})
					return
				}

				err := json.NewEncoder(w).Encode(tt.response)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)

					_ = json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("invalid response %s", err.Error())})
					return
				}
			})

			ts := httptest.NewServer(mockHandler)
			defer ts.Close()

			service := NewService(provider.NewClient(ts.URL))

			ctx := context.Background()

			result, err := service.List(ctx, tt.request)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %+v,\n got %+v", &tt.expected, result)
			}
		})
	}
}
