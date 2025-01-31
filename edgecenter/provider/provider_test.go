package provider

import (
	"context"
	"encoding/json"
	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
		name    string
		fields  fields
		args    args
		wantErr bool
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			tt.fields.baseURL = ts.URL

			c := &Client{
				httpc:   tt.fields.httpc,
				signer:  tt.fields.signer,
				ua:      tt.fields.ua,
				baseURL: tt.fields.baseURL,
			}

			err := c.Request(tt.args.ctx, tt.args.method, tt.args.path, tt.args.payload, tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("Request() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
