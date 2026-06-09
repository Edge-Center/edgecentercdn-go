package edgecentercdn_go

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Edge-Center/edgecentercdn-go/internal/testutil"
)

var _ ClientService = (*Service)(nil)

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
		{"Rules", func() interface{} { return svc.Rules() }},
		{"LECerts", func() interface{} { return svc.LECerts() }},
		{"OriginGroups", func() interface{} { return svc.OriginGroups() }},
		{"Shielding", func() interface{} { return svc.Shielding() }},
		{"SSLCerts", func() interface{} { return svc.SSLCerts() }},
		{"Statistics", func() interface{} { return svc.Statistics() }},
		{"Tools", func() interface{} { return svc.Tools() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.accessor())
		})
	}
}
