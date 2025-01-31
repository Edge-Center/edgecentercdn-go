package resources

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_buildListPath(t *testing.T) {
	tests := []struct {
		name     string
		req      *ListRequest
		expected string
	}{
		{
			name:     "Empty request",
			req:      &ListRequest{},
			expected: "/cdn/resources",
		},
		{
			name:     "Request with offset",
			req:      &ListRequest{Offset: 10},
			expected: "/cdn/resources?offset=10",
		},
		{
			name:     "Request with size",
			req:      &ListRequest{Size: 20},
			expected: "/cdn/resources?size=20",
		},
		{
			name:     "Request with fields",
			req:      &ListRequest{Filter: &ListFilterRequest{Fields: []string{"name", "status"}}},
			expected: "/cdn/resources?fields=name%2Cstatus",
		},
		{
			name:     "Request with ordering",
			req:      &ListRequest{Filter: &ListFilterRequest{Ordering: "desc"}},
			expected: "/cdn/resources?ordering=desc",
		},
		{
			name:     "Request with status",
			req:      &ListRequest{Filter: &ListFilterRequest{Status: []ResourceStatus{ActiveResourceStatus, ProcessedResourceStatus}}},
			expected: "/cdn/resources?status=active%2Cprocessed",
		},
		{
			name: "Request with multiple parameters",
			req: &ListRequest{
				Offset: 10, Size: 20, Filter: &ListFilterRequest{Fields: []string{"name", "status"}, Ordering: "asc", Status: []ResourceStatus{ActiveResourceStatus, ProcessedResourceStatus}},
			},
			expected: "/cdn/resources?fields=name%2Cstatus&offset=10&ordering=asc&size=20&status=active%2Cprocessed",
		},
		{
			name:     "Request with empty fields and status",
			req:      &ListRequest{Filter: &ListFilterRequest{Fields: []string{}, Status: []ResourceStatus{}}},
			expected: "/cdn/resources",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := buildListPath(tt.req.Offset, tt.req.Size, tt.req.Filter)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
