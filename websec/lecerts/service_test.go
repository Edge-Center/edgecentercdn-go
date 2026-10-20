package lecerts

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

func expectNoResultDecoding(t *testing.T, mock *testutil.MockRequester) {
	t.Helper()
	mock.RequestFunc = func(_ context.Context, _, _ string, _, result interface{}) error {
		assert.Nil(t, result)
		return nil
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

func TestService_Issue(t *testing.T) {
	tests := []struct {
		name        string
		req         *IssueRequest
		wantPayload string
	}{
		{name: "nil request sends no body", req: nil},
		{name: "empty cert type sends no body", req: &IssueRequest{}},
		{name: "MDDC", req: &IssueRequest{CertType: CertTypeMDDC}, wantPayload: `{"cert_type":"MDDC"}`},
		{name: "EC", req: &IssueRequest{CertType: CertTypeEC}, wantPayload: `{"cert_type":"EC"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &testutil.MockRequester{}
			expectNoResultDecoding(t, mock)

			err := NewService(mock).Issue(context.Background(), 100, tt.req)

			require.NoError(t, err)
			require.Len(t, mock.Calls, 1)
			assert.Equal(t, http.MethodPost, mock.Calls[0].Method)
			assert.Equal(t, "/resources/100/ssl/le/issue", mock.Calls[0].Path)
			if tt.wantPayload == "" {
				assert.Nil(t, mock.Calls[0].Payload)
				return
			}
			assert.JSONEq(t, tt.wantPayload, payloadJSON(t, mock.Calls[0].Payload))
		})
	}
}

func TestService_Renew(t *testing.T) {
	mock := &testutil.MockRequester{}
	expectNoResultDecoding(t, mock)

	err := NewService(mock).Renew(context.Background(), 100)

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodPost, mock.Calls[0].Method)
	assert.Equal(t, "/resources/100/ssl/le/renew", mock.Calls[0].Path)
	assert.Nil(t, mock.Calls[0].Payload)
}

func TestService_Revoke(t *testing.T) {
	tests := []struct {
		name     string
		force    bool
		wantPath string
	}{
		{name: "force", force: true, wantPath: "/resources/100/ssl/le/revoke?force=true"},
		{name: "no force", force: false, wantPath: "/resources/100/ssl/le/revoke?force=false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &testutil.MockRequester{}
			expectNoResultDecoding(t, mock)

			err := NewService(mock).Revoke(context.Background(), 100, tt.force)

			require.NoError(t, err)
			require.Len(t, mock.Calls, 1)
			assert.Equal(t, http.MethodPost, mock.Calls[0].Method)
			assert.Equal(t, tt.wantPath, mock.Calls[0].Path)
			assert.Nil(t, mock.Calls[0].Payload)
		})
	}
}

func TestService_Get(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage(`{
		"id": 42,
		"statuses": [
			{"id": 7, "status": "PROCESSING", "error": "", "details": "", "created": "2026-10-20T10:00:00.123456Z"},
			{"id": 8, "status": "DONE", "error": "", "details": "", "created": "2026-10-20T10:05:00.123456Z"}
		],
		"ssl_id": 555,
		"cert_type": "LE",
		"started": "2026-10-20T10:00:00.000001Z",
		"finished": "2026-10-20T10:05:00.000001Z",
		"acme_cert_type": "LE",
		"active": false,
		"resource": 100
	}`))

	status, err := NewService(mock).Get(context.Background(), 100)

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodGet, mock.Calls[0].Method)
	assert.Equal(t, "/resources/100/ssl/le/status", mock.Calls[0].Path)
	assert.Nil(t, mock.Calls[0].Payload)

	assert.Equal(t, &Status{
		ID:           42,
		Resource:     100,
		Active:       false,
		Started:      ptr("2026-10-20T10:00:00.000001Z"),
		Finished:     ptr("2026-10-20T10:05:00.000001Z"),
		ACMECertType: CertTypeLE,
		SSLID:        ptr[int64](555),
		CertType:     ptr(CertTypeLE),
		Statuses: []StatusDetail{
			{ID: 7, Status: "PROCESSING", Created: "2026-10-20T10:00:00.123456Z"},
			{ID: 8, Status: "DONE", Created: "2026-10-20T10:05:00.123456Z"},
		},
	}, status)
}

func TestService_Get_RunningWithoutAutomatedCertificate(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage(`{
		"id": 43,
		"statuses": [
			{"id": 9, "status": "FAILED", "error": "rateLimited", "details": "too many certificates", "created": "2026-10-20T10:15:00Z"}
		],
		"ssl_id": null,
		"cert_type": null,
		"started": null,
		"finished": null,
		"acme_cert_type": "MDDC",
		"active": true,
		"resource": 100
	}`))

	status, err := NewService(mock).Get(context.Background(), 100)

	require.NoError(t, err)
	assert.Equal(t, &Status{
		ID:           43,
		Resource:     100,
		Active:       true,
		ACMECertType: CertTypeMDDC,
		Statuses: []StatusDetail{
			{ID: 9, Status: "FAILED", Error: "rateLimited", Details: "too many certificates", Created: "2026-10-20T10:15:00Z"},
		},
	}, status)
}

func TestService_Cancel(t *testing.T) {
	mock := &testutil.MockRequester{}
	expectNoResultDecoding(t, mock)

	err := NewService(mock).Cancel(context.Background(), 100)

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodPut, mock.Calls[0].Method)
	assert.Equal(t, "/resources/100/ssl/le/status", mock.Calls[0].Path)
	assert.JSONEq(t, `{"active":false}`, payloadJSON(t, mock.Calls[0].Payload))
}

func TestService_ErrorPropagation(t *testing.T) {
	errRemote := errors.New("connection refused")

	tests := []struct {
		name string
		call func(svc Service) (interface{}, error)
	}{
		{"Issue", func(svc Service) (interface{}, error) {
			return nil, svc.Issue(context.Background(), 1, &IssueRequest{CertType: CertTypeLE})
		}},
		{"Renew", func(svc Service) (interface{}, error) { return nil, svc.Renew(context.Background(), 1) }},
		{"Revoke", func(svc Service) (interface{}, error) { return nil, svc.Revoke(context.Background(), 1, true) }},
		{"Get", func(svc Service) (interface{}, error) { return svc.Get(context.Background(), 1) }},
		{"Cancel", func(svc Service) (interface{}, error) { return nil, svc.Cancel(context.Background(), 1) }},
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
