package websec

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Edge-Center/edgecentercdn-go/edgecenter/provider"
	"github.com/Edge-Center/edgecentercdn-go/internal/testutil"
)

func TestNewService_NotNil(t *testing.T) {
	mock := &testutil.MockRequester{}
	svc := NewService(mock)
	require.NotNil(t, svc)
}

func TestService_Accessors(t *testing.T) {
	mock := &testutil.MockRequester{}
	svc := NewService(mock)

	tests := []struct {
		name     string
		accessor func() interface{}
	}{
		{"Resources", func() interface{} { return svc.Resources() }},
		{"Locations", func() interface{} { return svc.Locations() }},
		{"LECerts", func() interface{} { return svc.LECerts() }},
		{"SourceGroups", func() interface{} { return svc.SourceGroups() }},
		{"SSLCerts", func() interface{} { return svc.SSLCerts() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.accessor())
		})
	}
}

func TestService_SendsRelativePathsUnderTheBaseURL(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		status     int
		body       string
		call       func(svc *Service) error
		wantMethod string
		wantURI    string
	}{
		{
			name:   "le revoke",
			status: http.StatusOK,
			body:   `{"message": "Certificate has been revoked"}`,
			call: func(svc *Service) error {
				return svc.LECerts().Revoke(ctx, 7, true)
			},
			wantMethod: http.MethodPost,
			wantURI:    "/protection/resources/7/ssl/le/revoke?force=true",
		},
		{
			name:   "source group delete",
			status: http.StatusNoContent,
			call: func(svc *Service) error {
				return svc.SourceGroups().Delete(ctx, 5)
			},
			wantMethod: http.MethodDelete,
			wantURI:    "/protection/source_groups/5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var method, uri string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				method = r.Method
				uri = r.URL.RequestURI()
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			t.Cleanup(srv.Close)
			svc := NewService(provider.NewClient(srv.URL + "/protection/"))

			require.NoError(t, tt.call(svc))
			assert.Equal(t, tt.wantMethod, method)
			assert.Equal(t, tt.wantURI, uri)
		})
	}
}
