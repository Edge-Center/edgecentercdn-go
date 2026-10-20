package sslcerts

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
)

const certResponseJSON = `{
	"id": 7,
	"name": "example",
	"deleted": false,
	"cert_issuer": "R3",
	"cert_subject_cn": "example.com",
	"cert_subject_alt": "example.com,www.example.com",
	"validity_not_before": "2026-01-01T00:00:00Z",
	"validity_not_after": "2026-04-01T12:30:00.123456Z",
	"certificates_chain": "gAAAAABencryptedchain",
	"related_resources": [{"id": 11, "domain": "example.com"}],
	"has_related_resources": true,
	"automated": true
}`

func expectedCert() Cert {
	notBefore := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	notAfter := time.Date(2026, 4, 1, 12, 30, 0, 123456000, time.UTC)

	return Cert{
		ID:                  7,
		Name:                "example",
		Automated:           true,
		CertIssuer:          "R3",
		CertSubjectCN:       "example.com",
		CertSubjectAlt:      "example.com,www.example.com",
		ValidityNotBefore:   &notBefore,
		ValidityNotAfter:    &notAfter,
		RelatedResources:    []RelatedResource{{ID: 11, Domain: "example.com"}},
		HasRelatedResources: true,
	}
}

func validCreateRequest() *CreateRequest {
	return &CreateRequest{
		Name:        "example",
		Certificate: "-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----\n",
		PrivateKey:  "-----BEGIN PRIVATE KEY-----\nMIIE\n-----END PRIVATE KEY-----\n",
	}
}

func payloadJSON(t *testing.T, payload interface{}) string {
	t.Helper()
	b, err := json.Marshal(payload)
	require.NoError(t, err)
	return string(b)
}

func TestService_List(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage("[" + certResponseJSON + `,{"id": 8, "name": "second"}]`))

	svc := NewService(mock)
	certs, err := svc.List(context.Background())

	require.NoError(t, err)
	require.Len(t, certs, 2)
	assert.Equal(t, expectedCert(), certs[0])
	assert.Equal(t, Cert{ID: 8, Name: "second"}, certs[1])
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodGet, mock.Calls[0].Method)
	assert.Equal(t, "/ssl/certificates", mock.Calls[0].Path)
	assert.Nil(t, mock.Calls[0].Payload)
}

func TestService_Get(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage(certResponseJSON))

	svc := NewService(mock)
	cert, err := svc.Get(context.Background(), 7)

	require.NoError(t, err)
	assert.Equal(t, expectedCert(), *cert)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodGet, mock.Calls[0].Method)
	assert.Equal(t, "/ssl/certificates/7", mock.Calls[0].Path)
	assert.Nil(t, mock.Calls[0].Payload)
}

func TestService_Create(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage(certResponseJSON))

	req := validCreateRequest()
	req.CertificatesChain = "-----BEGIN CERTIFICATE-----\nMIIC\n-----END CERTIFICATE-----\n"

	svc := NewService(mock)
	cert, err := svc.Create(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, expectedCert(), *cert)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodPost, mock.Calls[0].Method)
	assert.Equal(t, "/ssl/certificates", mock.Calls[0].Path)
	assert.JSONEq(t, `{
		"name": "example",
		"certificate": "-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----\n",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIE\n-----END PRIVATE KEY-----\n",
		"certificates_chain": "-----BEGIN CERTIFICATE-----\nMIIC\n-----END CERTIFICATE-----\n"
	}`, payloadJSON(t, mock.Calls[0].Payload))
}

func TestService_Create_OmitsEmptyChain(t *testing.T) {
	mock := &testutil.MockRequester{}

	svc := NewService(mock)
	_, err := svc.Create(context.Background(), validCreateRequest())

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.JSONEq(t, `{
		"name": "example",
		"certificate": "-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----\n",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIE\n-----END PRIVATE KEY-----\n"
	}`, payloadJSON(t, mock.Calls[0].Payload))
}

func TestService_Create_ValidateError(t *testing.T) {
	mock := &testutil.MockRequester{}
	req := validCreateRequest()
	req.PrivateKey = ""

	svc := NewService(mock)
	cert, err := svc.Create(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, cert)
	assert.Equal(t, "validate ssl cert create request: private_key is required", err.Error())
	assert.Empty(t, mock.Calls)
}

func TestService_Update_SendsOnlyName(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(json.RawMessage(certResponseJSON))

	svc := NewService(mock)
	cert, err := svc.Update(context.Background(), 7, &UpdateRequest{Name: "renamed"})

	require.NoError(t, err)
	assert.Equal(t, expectedCert(), *cert)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodPut, mock.Calls[0].Method)
	assert.Equal(t, "/ssl/certificates/7", mock.Calls[0].Path)
	assert.JSONEq(t, `{"name": "renamed"}`, payloadJSON(t, mock.Calls[0].Payload))
}

func TestService_Update_ValidateError(t *testing.T) {
	mock := &testutil.MockRequester{}

	svc := NewService(mock)
	cert, err := svc.Update(context.Background(), 7, &UpdateRequest{})

	require.Error(t, err)
	assert.Nil(t, cert)
	assert.Equal(t, "validate ssl cert update request: name is required", err.Error())
	assert.Empty(t, mock.Calls)
}

func TestService_Delete(t *testing.T) {
	mock := &testutil.MockRequester{}

	svc := NewService(mock)
	err := svc.Delete(context.Background(), 7)

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodDelete, mock.Calls[0].Method)
	assert.Equal(t, "/ssl/certificates/7", mock.Calls[0].Path)
	assert.Nil(t, mock.Calls[0].Payload)
}

func TestService_ErrorPropagation(t *testing.T) {
	errRemote := errors.New("connection refused")

	tests := []struct {
		name string
		call func(svc Service) (interface{}, error)
	}{
		{"List", func(svc Service) (interface{}, error) { return svc.List(context.Background()) }},
		{"Get", func(svc Service) (interface{}, error) { return svc.Get(context.Background(), 7) }},
		{"Create", func(svc Service) (interface{}, error) { return svc.Create(context.Background(), validCreateRequest()) }},
		{"Update", func(svc Service) (interface{}, error) {
			return svc.Update(context.Background(), 7, &UpdateRequest{Name: "renamed"})
		}},
		{"Delete", func(svc Service) (interface{}, error) { return nil, svc.Delete(context.Background(), 7) }},
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
			req:  &CreateRequest{Name: "example", Certificate: "cert", PrivateKey: "key"},
		},
		{
			name:    "missing name",
			req:     &CreateRequest{Certificate: "cert", PrivateKey: "key"},
			wantErr: "name is required",
		},
		{
			name:    "missing certificate",
			req:     &CreateRequest{Name: "example", PrivateKey: "key"},
			wantErr: "certificate is required",
		},
		{
			name:    "missing private key",
			req:     &CreateRequest{Name: "example", Certificate: "cert"},
			wantErr: "private_key is required",
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
