package statistics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_buildStatisticsPath(t *testing.T) {
	fromTime := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)
	toTime := time.Date(2023, 1, 15, 1, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		request      *ResourceStatisticsTimeSeriesRequest
		expectedPath string
	}{
		{
			name: "All parameters present",
			request: &ResourceStatisticsTimeSeriesRequest{
				Metrics:     []Metric{MetricResponses2xx, MetricTotalBytes},
				Regions:     []Region{RegionNA, RegionCIS},
				GroupBy:     []GroupBy{GroupByRegion, GroupByResource},
				Hosts:       []Host{"cdn.example.com", "cdn1.example.com"},
				Resources:   []ResourceId{1, 2},
				Granularity: Granularity1h,
				From:        fromTime,
				To:          toTime,
			},
			expectedPath: "/cdn/statistics/aggregate/stats?flat=true&from=2023-01-15T00%3A00%3A00Z&granularity=1h&group_by=region&group_by=resource&metrics=responses_2xx&metrics=total_bytes&regions=na&regions=cis&resource=1&resource=2&service=CDN&to=2023-01-15T01%3A00%3A00Z&vhost=cdn.example.com&vhost=cdn1.example.com",
		},
		{
			name: "Group by regions",
			request: &ResourceStatisticsTimeSeriesRequest{
				Metrics:     []Metric{MetricResponses2xx},
				Regions:     []Region{RegionNA},
				GroupBy:     []GroupBy{GroupByRegion},
				Granularity: Granularity1h,
				From:        fromTime,
				To:          toTime,
			},
			expectedPath: "/cdn/statistics/aggregate/stats?flat=true&from=2023-01-15T00%3A00%3A00Z&granularity=1h&group_by=region&metrics=responses_2xx&regions=na&service=CDN&to=2023-01-15T01%3A00%3A00Z",
		},
		{
			name: "Group by countries",
			request: &ResourceStatisticsTimeSeriesRequest{
				Metrics:     []Metric{MetricResponses2xx},
				Regions:     []Region{RegionNA},
				GroupBy:     []GroupBy{GroupByCountry},
				Countries:   []Country{"Russian Federation"},
				Granularity: Granularity1h,
				From:        fromTime,
				To:          toTime,
			},
			expectedPath: "/cdn/statistics/aggregate/stats?countries=Russian+Federation&flat=true&from=2023-01-15T00%3A00%3A00Z&granularity=1h&group_by=country&metrics=responses_2xx&regions=na&service=CDN&to=2023-01-15T01%3A00%3A00Z",
		},
		{
			name: "Group by resources",
			request: &ResourceStatisticsTimeSeriesRequest{
				Metrics:     []Metric{MetricResponses2xx},
				Resources:   []ResourceId{1, 2},
				GroupBy:     []GroupBy{GroupByResource},
				Granularity: Granularity1h,
				From:        fromTime,
				To:          toTime,
			},
			expectedPath: "/cdn/statistics/aggregate/stats?flat=true&from=2023-01-15T00%3A00%3A00Z&granularity=1h&group_by=resource&metrics=responses_2xx&resource=1&resource=2&service=CDN&to=2023-01-15T01%3A00%3A00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildStatisticsPath(
				tt.request.From,
				tt.request.To,
				tt.request.Granularity,
				tt.request.GroupBy,
				tt.request.Metrics,
				tt.request.Regions,
				tt.request.Hosts,
				tt.request.Resources,
				tt.request.Clients,
				tt.request.Countries,
			)

			assert.Equal(t, tt.expectedPath, result)
		})
	}
}

func TestParseTimeRange(t *testing.T) {
	tests := []struct {
		name      string
		from      string
		to        string
		expectErr bool
	}{
		{
			name:      "Valid from and to",
			from:      "2025-01-01T12:00:00Z",
			to:        "2025-01-02T12:00:00Z",
			expectErr: false,
		},
		{
			name:      "Invalid from format",
			from:      "invalid",
			to:        "2025-01-02T12:00:00Z",
			expectErr: true,
		},
		{
			name:      "Invalid to format",
			from:      "2025-01-01T12:00:00Z",
			to:        "invalid",
			expectErr: true,
		},
		{
			name:      "Empty from and to",
			from:      "",
			to:        "",
			expectErr: false,
		},
		{
			name:      "Empty from",
			from:      "",
			to:        "2025-01-02T12:00:00Z",
			expectErr: false,
		},
		{
			name:      "Empty to",
			from:      "2025-01-01T12:00:00Z",
			to:        "",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fromTime, toTime, err := ParseTimeRange(tt.from, tt.to)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("did not expect an error but got: %v", err)
				}
				if tt.from != "" {
					expectedFrom, _ := time.Parse(time.RFC3339, tt.from)
					if !fromTime.Equal(expectedFrom) {
						t.Errorf("expected fromTime %v, got %v", expectedFrom, fromTime)
					}
				}
				if tt.to != "" {
					expectedTo, _ := time.Parse(time.RFC3339, tt.to)
					if !toTime.Equal(expectedTo) {
						t.Errorf("expected toTime %v, got %v", expectedTo, toTime)
					}
				}
			}
		})
	}
}

func TestGranularityToSeconds(t *testing.T) {
	tests := []struct {
		name          string
		granularity   Granularity
		expected      int64
		expectedError bool
	}{
		{
			name:        "Granularity 1m",
			granularity: Granularity1m,
			expected:    60,
		},
		{
			name:        "Granularity 5m",
			granularity: Granularity5m,
			expected:    300,
		},
		{
			name:        "Granularity 15m",
			granularity: Granularity15m,
			expected:    900,
		},
		{
			name:        "Granularity 1h",
			granularity: Granularity1h,
			expected:    3600,
		},
		{
			name:        "Granularity 1d",
			granularity: Granularity1d,
			expected:    86400,
		},
		{
			name:          "Unknown Granularity",
			granularity:   "2h",
			expected:      0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seconds, err := GranularityToSeconds(tt.granularity)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, seconds)
			}
		})
	}
}
