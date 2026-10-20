package resources

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Edge-Center/edgecentercdn-go/internal/testutil"
	"github.com/Edge-Center/edgecentercdn-go/websec/locations"
	"github.com/Edge-Center/edgecentercdn-go/websec/options"
)

const resourceResponseJSON = `{
	"id": 151464,
	"subdomains": ["www.example.com", "api.example.com"],
	"protected_ip": "198.51.100.7",
	"locations": [
		{
			"id": 2,
			"name": "info",
			"location": "^/info",
			"active": true,
			"source_group": null,
			"weight": 2,
			"deleted": false,
			"options": {"allowed_http_methods": {"enabled": true, "value": ["POST", "GET"]}}
		}
	],
	"client": 101,
	"options": {
		"blacklist": {"enabled": true, "value": ["150.150.150.150/32"]},
		"whitelist": null
	},
	"deleted": false,
	"active": true,
	"enabled": true,
	"domain": "example.com",
	"status": "active",
	"ssl_certificate": 4803,
	"ssl_automated": false,
	"source_group": 12,
	"created": "2026-04-03T10:40:59.477231Z",
	"updated": "2026-09-26T06:30:48.655789Z"
}`

func ptr[T any](v T) *T {
	return &v
}

func expectedResource() Resource {
	return Resource{
		ID:             151464,
		Domain:         "example.com",
		Subdomains:     []string{"www.example.com", "api.example.com"},
		Active:         true,
		Enabled:        true,
		Status:         StatusActive,
		ProtectedIP:    ptr("198.51.100.7"),
		SSLCertificate: ptr[int64](4803),
		SSLAutomated:   false,
		SourceGroup:    ptr[int64](12),
		Client:         101,
		Locations: []locations.Location{
			{
				ID:       2,
				Name:     "info",
				Location: "^/info",
				Active:   true,
				Weight:   2,
				Options: &options.LocationOptions{
					AllowedHTTPMethods: options.Value(options.ListOption{Enabled: true, Value: []string{"POST", "GET"}}),
				},
			},
		},
		Options: &options.ResourceOptions{
			Blacklist: options.Value(options.ListOption{Enabled: true, Value: []string{"150.150.150.150/32"}}),
			Whitelist: options.Null[options.ListOption](),
		},
		Created: time.Date(2026, 4, 3, 10, 40, 59, 477231000, time.UTC),
		Updated: time.Date(2026, 9, 26, 6, 30, 48, 655789000, time.UTC),
	}
}

func payloadJSON(t *testing.T, payload interface{}) string {
	t.Helper()
	b, err := json.Marshal(payload)
	require.NoError(t, err)
	return string(b)
}

func validUpdateRequest() *UpdateRequest {
	return &UpdateRequest{
		Subdomains:  options.Value([]string{"www.example.com"}),
		Active:      true,
		SourceGroup: 12,
	}
}

func TestService_List(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage("[" + resourceResponseJSON + "]"))

	svc := NewService(mock)
	resources, err := svc.List(context.Background())

	require.NoError(t, err)
	assert.Equal(t, []Resource{expectedResource()}, resources)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodGet, mock.Calls[0].Method)
	assert.Equal(t, "/resources", mock.Calls[0].Path)
	assert.Nil(t, mock.Calls[0].Payload)
}

func TestService_Get(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage(resourceResponseJSON))

	svc := NewService(mock)
	resource, err := svc.Get(context.Background(), 151464)

	require.NoError(t, err)
	assert.Equal(t, expectedResource(), *resource)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodGet, mock.Calls[0].Method)
	assert.Equal(t, "/resources/151464", mock.Calls[0].Path)
	assert.Nil(t, mock.Calls[0].Payload)
}

func TestService_Create(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage(resourceResponseJSON))

	svc := NewService(mock)
	resource, err := svc.Create(context.Background(), &CreateRequest{
		Domain:         "example.com",
		Subdomains:     []string{"www.example.com", "api.example.com"},
		Active:         ptr(false),
		SSLCertificate: ptr[int64](4803),
		SSLAutomated:   true,
		SourceGroup:    12,
		Options: &options.ResourceOptions{
			RedirectHTTPToHTTPS: options.Value(options.BooleanOption{Enabled: true, Value: true}),
			Blacklist:           options.Value(options.ListOption{Enabled: true, Value: []string{"150.150.150.150"}}),
			WAF:                 options.Null[options.BooleanOption](),
		},
	})

	require.NoError(t, err)
	assert.Equal(t, expectedResource(), *resource)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodPost, mock.Calls[0].Method)
	assert.Equal(t, "/resources", mock.Calls[0].Path)
	assert.JSONEq(t, `{
		"domain": "example.com",
		"subdomains": ["www.example.com", "api.example.com"],
		"active": false,
		"ssl_certificate": 4803,
		"ssl_automated": true,
		"source_group": 12,
		"options": {
			"redirect_http_to_https": {"enabled": true, "value": true},
			"blacklist": {"enabled": true, "value": ["150.150.150.150"]},
			"waf": null
		}
	}`, payloadJSON(t, mock.Calls[0].Payload))
}

func TestService_Create_OmitsUnsetFieldsSoServerDefaultsApply(t *testing.T) {
	mock := &testutil.MockRequester{}

	svc := NewService(mock)
	_, err := svc.Create(context.Background(), &CreateRequest{
		Domain:      "example.org",
		SourceGroup: 12,
	})

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.JSONEq(t, `{"domain": "example.org", "source_group": 12}`, payloadJSON(t, mock.Calls[0].Payload))
}

func TestService_Create_ValidateError(t *testing.T) {
	mock := &testutil.MockRequester{}

	svc := NewService(mock)
	resource, err := svc.Create(context.Background(), &CreateRequest{SourceGroup: 12})

	require.Error(t, err)
	assert.Nil(t, resource)
	assert.Equal(t, "validate resource create request: domain is required", err.Error())
	assert.Empty(t, mock.Calls)
}

func TestService_Update(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage(resourceResponseJSON))

	svc := NewService(mock)
	resource, err := svc.Update(context.Background(), 151464, &UpdateRequest{
		Subdomains:     options.Value([]string{"www.example.com", "api.example.com"}),
		Active:         true,
		SSLCertificate: options.Value[int64](4803),
		SSLAutomated:   ptr(false),
		SourceGroup:    12,
		Options: &options.ResourceOptions{
			TLSVersions: options.Value(options.ListOption{Enabled: true, Value: []string{"TLSv1.2", "TLSv1.3"}}),
			Whitelist:   options.Null[options.ListOption](),
		},
	})

	require.NoError(t, err)
	assert.Equal(t, expectedResource(), *resource)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodPut, mock.Calls[0].Method)
	assert.Equal(t, "/resources/151464", mock.Calls[0].Path)
	assert.JSONEq(t, `{
		"subdomains": ["www.example.com", "api.example.com"],
		"active": true,
		"ssl_certificate": 4803,
		"ssl_automated": false,
		"source_group": 12,
		"options": {
			"tls_versions": {"enabled": true, "value": ["TLSv1.2", "TLSv1.3"]},
			"whitelist": null
		}
	}`, payloadJSON(t, mock.Calls[0].Payload))
}

func TestService_Update_SSLCertificateStates(t *testing.T) {
	tests := []struct {
		name string
		cert options.Field[int64]
		want string
	}{
		{
			name: "absent keeps the linked certificate",
			cert: options.Field[int64]{},
			want: `{"subdomains": ["www.example.com"], "active": true, "source_group": 12}`,
		},
		{
			name: "null unlinks the certificate",
			cert: options.Null[int64](),
			want: `{"subdomains": ["www.example.com"], "active": true, "ssl_certificate": null, "source_group": 12}`,
		},
		{
			name: "value links the certificate",
			cert: options.Value[int64](4803),
			want: `{"subdomains": ["www.example.com"], "active": true, "ssl_certificate": 4803, "source_group": 12}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &testutil.MockRequester{}
			req := validUpdateRequest()
			req.SSLCertificate = tt.cert

			svc := NewService(mock)
			_, err := svc.Update(context.Background(), 151464, req)

			require.NoError(t, err)
			require.Len(t, mock.Calls, 1)
			assert.JSONEq(t, tt.want, payloadJSON(t, mock.Calls[0].Payload))
		})
	}
}

func TestService_Update_SubdomainsStates(t *testing.T) {
	tests := []struct {
		name       string
		subdomains options.Field[[]string]
		want       string
	}{
		{
			name:       "absent leaves the subdomains alone",
			subdomains: options.Field[[]string]{},
			want:       `{"active": true, "source_group": 12}`,
		},
		{
			name:       "empty value clears the subdomains",
			subdomains: options.Value([]string{}),
			want:       `{"subdomains": [], "active": true, "source_group": 12}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &testutil.MockRequester{}
			req := validUpdateRequest()
			req.Subdomains = tt.subdomains

			svc := NewService(mock)
			_, err := svc.Update(context.Background(), 151464, req)

			require.NoError(t, err)
			require.Len(t, mock.Calls, 1)
			assert.JSONEq(t, tt.want, payloadJSON(t, mock.Calls[0].Payload))
		})
	}
}

func TestService_Update_ValidateError(t *testing.T) {
	mock := &testutil.MockRequester{}
	req := validUpdateRequest()
	req.SourceGroup = 0

	svc := NewService(mock)
	resource, err := svc.Update(context.Background(), 151464, req)

	require.Error(t, err)
	assert.Nil(t, resource)
	assert.Equal(t, "validate resource update request: source_group is required", err.Error())
	assert.Empty(t, mock.Calls)
}

func TestService_Delete(t *testing.T) {
	mock := &testutil.MockRequester{}

	svc := NewService(mock)
	err := svc.Delete(context.Background(), 151464)

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodDelete, mock.Calls[0].Method)
	assert.Equal(t, "/resources/151464", mock.Calls[0].Path)
	assert.Nil(t, mock.Calls[0].Payload)
}

func TestService_ErrorPropagation(t *testing.T) {
	errRemote := errors.New("connection refused")

	tests := []struct {
		name string
		call func(svc Service) (interface{}, error)
	}{
		{"List", func(svc Service) (interface{}, error) { return svc.List(context.Background()) }},
		{"Get", func(svc Service) (interface{}, error) { return svc.Get(context.Background(), 1) }},
		{"Create", func(svc Service) (interface{}, error) {
			return svc.Create(context.Background(), &CreateRequest{Domain: "example.com", SourceGroup: 12})
		}},
		{"Update", func(svc Service) (interface{}, error) {
			return svc.Update(context.Background(), 1, validUpdateRequest())
		}},
		{"Delete", func(svc Service) (interface{}, error) { return nil, svc.Delete(context.Background(), 1) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &testutil.MockRequester{}
			mock.RespondWithError(errRemote)

			result, err := tt.call(NewService(mock))

			require.Error(t, err)
			assert.Nil(t, result)
			assert.ErrorIs(t, err, errRemote)
			assert.Equal(t, "request: connection refused", err.Error())
			require.Len(t, mock.Calls, 1)
		})
	}
}

func TestCreateRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateRequest
		wantErr string
	}{
		{
			name: "only required fields",
			req:  &CreateRequest{Domain: "example.com", SourceGroup: 12},
		},
		{
			name:    "missing domain",
			req:     &CreateRequest{SourceGroup: 12},
			wantErr: "domain is required",
		},
		{
			name:    "missing source group",
			req:     &CreateRequest{Domain: "example.com"},
			wantErr: "source_group is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()

			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			assert.Equal(t, tt.wantErr, err.Error())
		})
	}
}
