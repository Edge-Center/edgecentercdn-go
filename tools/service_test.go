package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Edge-Center/edgecentercdn-go/edgecenter/provider"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestService_Purge(t *testing.T) {
	tests := []struct {
		name                string
		resourceId          int64
		request             *PurgeRequest
		response            PurgeResponse
		expectedResponse    PurgeResponse
		expectedRequestBody string
	}{
		{
			name:       "Purge paths",
			resourceId: 1,
			request: &PurgeRequest{
				Paths: []string{"/en.json"},
			},
			response: PurgeResponse{
				Paths: []string{"/en.json"},
			},
			expectedResponse: PurgeResponse{
				Paths: []string{"/en.json"},
			},
			expectedRequestBody: `{"paths":["/en.json"]}`,
		},
		{
			name:       "Empty paths",
			resourceId: 1,
			request:    nil,
			response: PurgeResponse{
				Paths: []string{},
			},
			expectedResponse: PurgeResponse{
				Paths: []string{},
			},
			expectedRequestBody: `{"paths":[]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != fmt.Sprintf("/cdn/resources/%d/purge", tt.resourceId) || r.Method != http.MethodPost {
					http.Error(w, "not found", http.StatusNotFound)
					return
				}

				w.Header().Set("Content-Type", "application/json")

				reqBody, err := io.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(map[string]string{"message": "could not read request body"})
					return
				}

				if strings.TrimSpace(string(reqBody)) != tt.expectedRequestBody {
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(map[string]string{"message": "invalid request body"})
					return
				}

				if err := json.NewEncoder(w).Encode(tt.response); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					_ = json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("invalid response %s", err.Error())})
					return
				}
			})

			ts := httptest.NewServer(mockHandler)
			defer ts.Close()

			service := NewService(provider.NewClient(ts.URL))

			ctx := context.Background()

			response, err := service.Purge(ctx, tt.resourceId, tt.request)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(response, &tt.expectedResponse) {
				t.Errorf("expected %+v, got %+v", &tt.expectedResponse, response)
			}
		})
	}
}
