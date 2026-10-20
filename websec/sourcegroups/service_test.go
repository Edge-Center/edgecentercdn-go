package sourcegroups

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Edge-Center/edgecentercdn-go/internal/testutil"
)

const groupResponseJSON = `{
	"id": 42,
	"name": "origins",
	"use_next": true,
	"description": "main origins",
	"sources": [
		{
			"comment": "main",
			"fail_timeout": 10,
			"max_fails": 3,
			"mode": "primary",
			"source": "203.0.113.10:8080",
			"weight": 60
		},
		{
			"comment": "",
			"fail_timeout": 1,
			"max_fails": 1,
			"mode": "backup",
			"source": "backup.example.com",
			"weight": 40
		}
	],
	"related_resources": [
		{"id": 11, "domain": "example.com", "locations": [{"id": 5, "name": "static"}]},
		{"id": 12, "domain": "other.example.com", "locations": []}
	]
}`

func expectedGroup() SourceGroup {
	return SourceGroup{
		ID:          42,
		Name:        "origins",
		UseNext:     true,
		Description: "main origins",
		Sources: []Source{
			{
				Source:      "203.0.113.10:8080",
				Mode:        SourceModePrimary,
				Weight:      60,
				MaxFails:    3,
				FailTimeout: 10,
				Comment:     "main",
			},
			{
				Source:      "backup.example.com",
				Mode:        SourceModeBackup,
				Weight:      40,
				MaxFails:    1,
				FailTimeout: 1,
			},
		},
		RelatedResources: []RelatedResource{
			{ID: 11, Domain: "example.com", Locations: []RelatedLocation{{ID: 5, Name: "static"}}},
			{ID: 12, Domain: "other.example.com", Locations: []RelatedLocation{}},
		},
	}
}

func fullRequest() *GroupRequest {
	return &GroupRequest{
		Name:        "origins",
		UseNext:     ptr(false),
		Description: "main origins",
		Sources: []Source{
			{
				Source:      "203.0.113.10:8080",
				Mode:        SourceModePrimary,
				Weight:      60,
				MaxFails:    3,
				FailTimeout: 10,
				Comment:     "main",
			},
			{
				Source: "backup.example.com",
				Mode:   SourceModeDown,
				Weight: 40,
			},
		},
	}
}

func ptr[T any](v T) *T {
	return &v
}

func payloadJSON(t *testing.T, payload interface{}) string {
	t.Helper()
	b, err := json.Marshal(payload)
	require.NoError(t, err)
	return string(b)
}

func TestService_List(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage("[" + groupResponseJSON + "]"))

	svc := NewService(mock)
	groups, err := svc.List(context.Background())

	require.NoError(t, err)
	assert.Equal(t, []SourceGroup{expectedGroup()}, groups)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodGet, mock.Calls[0].Method)
	assert.Equal(t, "/source_groups", mock.Calls[0].Path)
	assert.Nil(t, mock.Calls[0].Payload)
}

func TestService_Get(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage(groupResponseJSON))

	svc := NewService(mock)
	group, err := svc.Get(context.Background(), 42)

	require.NoError(t, err)
	assert.Equal(t, expectedGroup(), *group)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodGet, mock.Calls[0].Method)
	assert.Equal(t, "/source_groups/42", mock.Calls[0].Path)
	assert.Nil(t, mock.Calls[0].Payload)
}

func TestService_Create(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage(groupResponseJSON))

	svc := NewService(mock)
	group, err := svc.Create(context.Background(), fullRequest())

	require.NoError(t, err)
	assert.Equal(t, expectedGroup(), *group)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodPost, mock.Calls[0].Method)
	assert.Equal(t, "/source_groups", mock.Calls[0].Path)
	assert.JSONEq(t, `{
		"name": "origins",
		"use_next": false,
		"description": "main origins",
		"sources": [
			{
				"source": "203.0.113.10:8080",
				"mode": "primary",
				"weight": 60,
				"max_fails": 3,
				"fail_timeout": 10,
				"comment": "main"
			},
			{
				"source": "backup.example.com",
				"mode": "down",
				"weight": 40
			}
		]
	}`, payloadJSON(t, mock.Calls[0].Payload))
}

func TestService_Create_LeavesUseNextAndSourceDefaultsToServer(t *testing.T) {
	mock := &testutil.MockRequester{}

	svc := NewService(mock)
	_, err := svc.Create(context.Background(), &GroupRequest{
		Name:    "minimal",
		Sources: []Source{{Source: "203.0.113.10"}},
	})

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.JSONEq(t, `{
		"name": "minimal",
		"description": "",
		"sources": [{"source": "203.0.113.10"}]
	}`, payloadJSON(t, mock.Calls[0].Payload))
}

func TestService_Create_ValidateError(t *testing.T) {
	mock := &testutil.MockRequester{}
	req := fullRequest()
	req.Sources = nil

	svc := NewService(mock)
	group, err := svc.Create(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, group)
	assert.Equal(t, "validate source group create request: sources is required", err.Error())
	assert.Empty(t, mock.Calls)
}

func TestService_Update(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage(groupResponseJSON))

	svc := NewService(mock)
	group, err := svc.Update(context.Background(), 42, &GroupRequest{
		Name:    "origins",
		UseNext: ptr(true),
		Sources: []Source{{Source: "203.0.113.10:8080", Mode: SourceModePrimary}},
	})

	require.NoError(t, err)
	assert.Equal(t, expectedGroup(), *group)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodPut, mock.Calls[0].Method)
	assert.Equal(t, "/source_groups/42", mock.Calls[0].Path)
	assert.JSONEq(t, `{
		"name": "origins",
		"use_next": true,
		"description": "",
		"sources": [{"source": "203.0.113.10:8080", "mode": "primary"}]
	}`, payloadJSON(t, mock.Calls[0].Payload))
}

func TestService_Update_ValidateError(t *testing.T) {
	mock := &testutil.MockRequester{}
	req := fullRequest()
	req.Name = ""

	svc := NewService(mock)
	group, err := svc.Update(context.Background(), 42, req)

	require.Error(t, err)
	assert.Nil(t, group)
	assert.Equal(t, "validate source group update request: name is required", err.Error())
	assert.Empty(t, mock.Calls)
}

func TestService_Delete(t *testing.T) {
	mock := &testutil.MockRequester{}

	svc := NewService(mock)
	err := svc.Delete(context.Background(), 42)

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodDelete, mock.Calls[0].Method)
	assert.Equal(t, "/source_groups/42", mock.Calls[0].Path)
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
		{"Create", func(svc Service) (interface{}, error) { return svc.Create(context.Background(), fullRequest()) }},
		{"Update", func(svc Service) (interface{}, error) { return svc.Update(context.Background(), 1, fullRequest()) }},
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

func TestGroupRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *GroupRequest
		wantErr string
	}{
		{
			name: "only required fields",
			req:  &GroupRequest{Name: "minimal", Sources: []Source{{Source: "203.0.113.10"}}},
		},
		{
			name:    "missing name",
			req:     &GroupRequest{Sources: []Source{{Source: "203.0.113.10"}}},
			wantErr: "name is required",
		},
		{
			name:    "missing sources",
			req:     &GroupRequest{Name: "origins"},
			wantErr: "sources is required",
		},
		{
			name:    "empty sources",
			req:     &GroupRequest{Name: "origins", Sources: []Source{}},
			wantErr: "sources is required",
		},
		{
			name: "source without host",
			req: &GroupRequest{
				Name:    "origins",
				Sources: []Source{{Source: "203.0.113.10"}, {Mode: SourceModeBackup, Weight: 10}},
			},
			wantErr: "sources[1].source is required",
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
