package locations

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Edge-Center/edgecentercdn-go/internal/testutil"
	"github.com/Edge-Center/edgecentercdn-go/websec/options"
)

const locationResponseJSON = `{
	"id": 2,
	"name": "api",
	"location": "^/api/",
	"active": true,
	"source_group": 12,
	"weight": 5,
	"deleted": false,
	"options": {
		"allowed_http_methods": {"enabled": true, "value": ["GET", "POST"]},
		"asn_acl": null
	}
}`

func ptr[T any](v T) *T {
	return &v
}

func expectedLocation() Location {
	return Location{
		ID:          2,
		Name:        "api",
		Location:    "^/api/",
		Active:      true,
		SourceGroup: ptr[int64](12),
		Weight:      5,
		Options: &options.LocationOptions{
			AllowedHTTPMethods: options.Value(options.ListOption{Enabled: true, Value: []string{"GET", "POST"}}),
			ASNACL:             options.Null[options.ASNACLOption](),
		},
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
		Name:        "api",
		Location:    "^/api/",
		Active:      true,
		SourceGroup: ptr[int64](12),
	}
}

func TestService_List(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage("[" + locationResponseJSON + "]"))

	svc := NewService(mock)
	locations, err := svc.List(context.Background(), 151464)

	require.NoError(t, err)
	assert.Equal(t, []Location{expectedLocation()}, locations)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodGet, mock.Calls[0].Method)
	assert.Equal(t, "/resources/151464/locations", mock.Calls[0].Path)
	assert.Nil(t, mock.Calls[0].Payload)
}

func TestService_Get(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage(locationResponseJSON))

	svc := NewService(mock)
	location, err := svc.Get(context.Background(), 151464, 2)

	require.NoError(t, err)
	assert.Equal(t, expectedLocation(), *location)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodGet, mock.Calls[0].Method)
	assert.Equal(t, "/resources/151464/locations/2", mock.Calls[0].Path)
	assert.Nil(t, mock.Calls[0].Payload)
}

func TestService_Create(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage(locationResponseJSON))

	svc := NewService(mock)
	location, err := svc.Create(context.Background(), 151464, &CreateRequest{
		Name:        "api",
		Location:    "^/api/",
		Active:      ptr(false),
		SourceGroup: ptr[int64](12),
		Weight:      ptr(5),
		Options: &options.LocationOptions{
			AllowedHTTPMethods: options.Value(options.ListOption{Enabled: true, Value: []string{"GET", "POST"}}),
			ASNACL: options.Value(options.ASNACLOption{
				Enabled:        true,
				PolicyType:     options.PolicyTypeDeny,
				ExceptedValues: []int64{64496, 64511},
			}),
			UserAgentACL: options.Null[options.AccessPolicyOption](),
		},
	})

	require.NoError(t, err)
	assert.Equal(t, expectedLocation(), *location)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodPost, mock.Calls[0].Method)
	assert.Equal(t, "/resources/151464/locations", mock.Calls[0].Path)
	assert.JSONEq(t, `{
		"name": "api",
		"location": "^/api/",
		"active": false,
		"source_group": 12,
		"weight": 5,
		"options": {
			"allowed_http_methods": {"enabled": true, "value": ["GET", "POST"]},
			"asn_acl": {"enabled": true, "policy_type": "deny", "excepted_values": [64496, 64511]},
			"user_agent_acl": null
		}
	}`, payloadJSON(t, mock.Calls[0].Payload))
}

func TestService_Create_OmitsUnsetFieldsSoServerDefaultsApply(t *testing.T) {
	mock := &testutil.MockRequester{}

	svc := NewService(mock)
	_, err := svc.Create(context.Background(), 151464, &CreateRequest{Location: "/"})

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.JSONEq(t, `{"location": "/"}`, payloadJSON(t, mock.Calls[0].Payload))
}

func TestService_Create_ValidateError(t *testing.T) {
	mock := &testutil.MockRequester{}

	svc := NewService(mock)
	location, err := svc.Create(context.Background(), 151464, &CreateRequest{Name: "api"})

	require.Error(t, err)
	assert.Nil(t, location)
	assert.Equal(t, "validate location create request: location is required", err.Error())
	assert.Empty(t, mock.Calls)
}

func TestService_Update(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage(locationResponseJSON))

	svc := NewService(mock)
	location, err := svc.Update(context.Background(), 151464, 2, &UpdateRequest{
		Name:        "api",
		Location:    "^/api/",
		Active:      true,
		SourceGroup: ptr[int64](12),
		Weight:      ptr(5),
		Options: &options.LocationOptions{
			AllowedHTTPMethods: options.Null[options.ListOption](),
		},
	})

	require.NoError(t, err)
	assert.Equal(t, expectedLocation(), *location)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodPut, mock.Calls[0].Method)
	assert.Equal(t, "/resources/151464/locations/2", mock.Calls[0].Path)
	assert.JSONEq(t, `{
		"name": "api",
		"location": "^/api/",
		"active": true,
		"source_group": 12,
		"weight": 5,
		"options": {"allowed_http_methods": null}
	}`, payloadJSON(t, mock.Calls[0].Payload))
}

func TestService_Update_AlwaysSendsActive(t *testing.T) {
	mock := &testutil.MockRequester{}

	svc := NewService(mock)
	_, err := svc.Update(context.Background(), 151464, 3, &UpdateRequest{Location: "/"})

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.JSONEq(t, `{"location": "/", "active": false, "source_group": null}`, payloadJSON(t, mock.Calls[0].Payload))
}

func TestService_Update_ValidateError(t *testing.T) {
	mock := &testutil.MockRequester{}
	req := validUpdateRequest()
	req.Location = ""

	svc := NewService(mock)
	location, err := svc.Update(context.Background(), 151464, 2, req)

	require.Error(t, err)
	assert.Nil(t, location)
	assert.Equal(t, "validate location update request: location is required", err.Error())
	assert.Empty(t, mock.Calls)
}

func TestService_Delete(t *testing.T) {
	mock := &testutil.MockRequester{}

	svc := NewService(mock)
	err := svc.Delete(context.Background(), 151464, 2)

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodDelete, mock.Calls[0].Method)
	assert.Equal(t, "/resources/151464/locations/2", mock.Calls[0].Path)
	assert.Nil(t, mock.Calls[0].Payload)
}

func TestService_ErrorPropagation(t *testing.T) {
	errRemote := errors.New("connection refused")

	tests := []struct {
		name string
		call func(svc Service) (interface{}, error)
	}{
		{"List", func(svc Service) (interface{}, error) { return svc.List(context.Background(), 1) }},
		{"Get", func(svc Service) (interface{}, error) { return svc.Get(context.Background(), 1, 2) }},
		{"Create", func(svc Service) (interface{}, error) {
			return svc.Create(context.Background(), 1, &CreateRequest{Location: "/"})
		}},
		{"Update", func(svc Service) (interface{}, error) {
			return svc.Update(context.Background(), 1, 2, validUpdateRequest())
		}},
		{"Delete", func(svc Service) (interface{}, error) { return nil, svc.Delete(context.Background(), 1, 2) }},
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
			name: "only location",
			req:  &CreateRequest{Location: "/"},
		},
		{
			name:    "missing location",
			req:     &CreateRequest{Name: "api", Weight: ptr(5)},
			wantErr: "location is required",
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
