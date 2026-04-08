package origingroups

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

func TestOriginGroupService_Create(t *testing.T) {
	expected := OriginGroup{
		ID:      1,
		Name:    "test-group",
		UseNext: true,
		Origins: []Origin{
			{ID: 1, Source: "origin.example.com", Enabled: true},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/source_groups", r.URL.Path)
		json.NewEncoder(w).Encode(expected)
	}))
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.Create(context.Background(), &GroupRequest{
		Name:    "test-group",
		UseNext: true,
		Origins: []OriginRequest{{Source: "origin.example.com", Enabled: true}},
	})

	require.NoError(t, err)
	assert.Equal(t, &expected, result)
}

func TestOriginGroupService_Get(t *testing.T) {
	expected := OriginGroup{
		ID:      1,
		Name:    "test-group",
		UseNext: true,
		Origins: []Origin{
			{ID: 1, Source: "origin.example.com", Enabled: true},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/cdn/source_groups/1", r.URL.Path)
		json.NewEncoder(w).Encode(expected)
	}))
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.Get(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, &expected, result)
}

func TestOriginGroupService_Update(t *testing.T) {
	group := OriginGroup{
		ID:   1,
		Name: "updated-group",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPut && r.URL.Path == "/cdn/source_groups/1":
			json.NewEncoder(w).Encode(group)
		case r.Method == http.MethodDelete && r.URL.Path == "/cdn/source_groups/1/authorization":
			w.WriteHeader(http.StatusOK)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.Update(context.Background(), 1, &GroupRequest{
		Name: "updated-group",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "updated-group", result.Name)
	assert.Nil(t, result.Authorization)
}

func TestOriginGroupService_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodDelete && r.URL.Path == "/cdn/source_groups/1":
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodDelete && r.URL.Path == "/cdn/source_groups/1/authorization":
			w.WriteHeader(http.StatusOK)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	err := service.Delete(context.Background(), 1)

	require.NoError(t, err)
}

func TestOriginGroupService_Create_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	_, err := service.Create(context.Background(), &GroupRequest{
		Name:    "test",
		Origins: []OriginRequest{{Source: "origin.example.com", Enabled: true}},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "request:")
}

func TestOriginGroupService_Get_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	_, err := service.Get(context.Background(), 999)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestOriginGroupService_ManageAuth_Create(t *testing.T) {
	group := OriginGroup{ID: 1, Name: "test-group"}
	authResp := Authorization{
		AuthType:    "aws_signature_v4",
		AccessKeyID: "AKID",
		BucketName:  "my-bucket",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/cdn/source_groups":
			json.NewEncoder(w).Encode(group)
		case r.Method == http.MethodPost && r.URL.Path == "/cdn/source_groups/1/authorization":
			json.NewEncoder(w).Encode(authResp)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.Create(context.Background(), &GroupRequest{
		Name:    "test-group",
		Origins: []OriginRequest{{Source: "origin.example.com", Enabled: true}},
		Authorization: &Authorization{
			AuthType:    "aws_signature_v4",
			AccessKeyID: "AKID",
			BucketName:  "my-bucket",
		},
	})

	require.NoError(t, err)
	require.NotNil(t, result.Authorization)
	assert.Equal(t, "aws_signature_v4", result.Authorization.AuthType)
	assert.Equal(t, "AKID", result.Authorization.AccessKeyID)
}

func TestOriginGroupService_ManageAuth_Update(t *testing.T) {
	existingAuth := Authorization{AuthType: "aws_signature_v4", AccessKeyID: "OLD"}
	group := OriginGroup{ID: 1, Name: "test-group", Authorization: &existingAuth}
	updatedAuth := Authorization{AuthType: "aws_signature_v4", AccessKeyID: "NEW"}

	var authMethod string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPut && r.URL.Path == "/cdn/source_groups/1":
			json.NewEncoder(w).Encode(group)
		case r.URL.Path == "/cdn/source_groups/1/authorization":
			authMethod = r.Method
			json.NewEncoder(w).Encode(updatedAuth)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.Update(context.Background(), 1, &GroupRequest{
		Name: "test-group",
		Authorization: &Authorization{
			AuthType:    "aws_signature_v4",
			AccessKeyID: "NEW",
		},
	})

	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, authMethod)
	require.NotNil(t, result.Authorization)
	assert.Equal(t, "NEW", result.Authorization.AccessKeyID)
}

func TestOriginGroupService_ManageAuth_Delete(t *testing.T) {
	group := OriginGroup{ID: 1, Name: "test-group", Authorization: &Authorization{AuthType: "aws_signature_v4"}}

	var authMethod string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPut && r.URL.Path == "/cdn/source_groups/1":
			json.NewEncoder(w).Encode(group)
		case r.URL.Path == "/cdn/source_groups/1/authorization":
			authMethod = r.Method
			w.WriteHeader(http.StatusOK)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.Update(context.Background(), 1, &GroupRequest{
		Name:          "test-group",
		Authorization: nil,
	})

	require.NoError(t, err)
	assert.Equal(t, http.MethodDelete, authMethod)
	assert.Nil(t, result.Authorization)
}

func TestOriginGroupService_ManageAuth_DeleteError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodDelete && r.URL.Path == "/cdn/source_groups/1":
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodDelete && r.URL.Path == "/cdn/source_groups/1/authorization":
			w.WriteHeader(http.StatusNotFound)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	err := service.Delete(context.Background(), 1)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestOriginGroupService_Create_ValidateError(t *testing.T) {
	service := NewService(provider.NewClient("http://example.com"))

	result, err := service.Create(context.Background(), &GroupRequest{
		Name: "test-group",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "validate origin group request: origins is required", err.Error())
}

func TestGroupRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *GroupRequest
		wantErr string
	}{
		{
			name: "valid request",
			req: &GroupRequest{
				Name:    "test-group",
				Origins: []OriginRequest{{Source: "origin.example.com", Enabled: true}},
			},
		},
		{
			name: "missing name",
			req: &GroupRequest{
				Origins: []OriginRequest{{Source: "origin.example.com", Enabled: true}},
			},
			wantErr: "name is required",
		},
		{
			name: "missing origins",
			req: &GroupRequest{
				Name: "test-group",
			},
			wantErr: "origins is required",
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
