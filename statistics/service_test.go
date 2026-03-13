package statistics

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Edge-Center/edgecentercdn-go/edgecenter/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceStatisticsService_GetTableData(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	resID := float64(100)

	expected := ResourceStatisticsTableResponse{
		{
			Resource: &resID,
			Metrics: struct {
				Requests              uint64                          `json:"requests,omitempty"`
				SentBytes             uint64                          `json:"sent_bytes,omitempty"`
				TotalBytes            uint64                          `json:"total_bytes,omitempty"`
				ShieldBytes           uint64                          `json:"shield_bytes,omitempty"`
				CDNBytes              uint64                          `json:"cdn_bytes,omitempty"`
				UpstreamBytes         uint64                          `json:"upstream_bytes,omitempty"`
				Responses2xx          uint64                          `json:"responses_2xx,omitempty"`
				Responses3xx          uint64                          `json:"responses_3xx,omitempty"`
				Responses4xx          uint64                          `json:"responses_4xx,omitempty"`
				Responses5xx          uint64                          `json:"responses_5xx,omitempty"`
				ResponsesHit          uint64                          `json:"responses_hit,omitempty"`
				ResponsesMiss         uint64                          `json:"responses_miss,omitempty"`
				ImageProcessed        uint64                          `json:"image_processed,omitempty"`
				OriginResponseTime    uint64                          `json:"origin_response_time,omitempty"`
				RequestTime           uint64                          `json:"request_time,omitempty"`
				RequestsWafPassed     uint64                          `json:"requests_waf_passed,omitempty"`
				CacheHitTrafficRatio  float64                         `json:"cache_hit_traffic_ratio,omitempty"`
				CacheHitRequestsRatio float64                         `json:"cache_hit_requests_ratio,omitempty"`
				ShieldTrafficRatio    float64                         `json:"shield_traffic_ratio,omitempty"`
				ShieldUsage           []ResourceStatisticsShieldUsage `json:"shield_usage,omitempty"`
			}{Requests: 1000, SentBytes: 5000},
		},
	}

	var capturedPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.RequestURI()
		assert.Equal(t, http.MethodGet, r.Method)
		json.NewEncoder(w).Encode(expected)
	}))
	defer ts.Close()

	req := &ResourceStatisticsTableRequest{
		Metrics:   []Metric{MetricRequests, MetricSentBytes},
		GroupBy:   []GroupBy{GroupByResource},
		Resources: []ResourceId{100},
		From:      from,
		To:        to,
	}

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.GetTableData(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, &expected, result)
	assert.Contains(t, capturedPath, "metrics=requests")
	assert.Contains(t, capturedPath, "metrics=sent_bytes")
	assert.Contains(t, capturedPath, "group_by=resource")
	assert.Contains(t, capturedPath, "resource=100")
}

func TestResourceStatisticsService_GetTimeSeriesData(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	region := "eu"

	expected := ResourceStatisticsTimeSeriesResponse{
		{
			Region: &region,
			Metrics: struct {
				Requests              [][]uint64                      `json:"requests,omitempty"`
				SentBytes             [][]uint64                      `json:"sent_bytes,omitempty"`
				TotalBytes            [][]uint64                      `json:"total_bytes,omitempty"`
				ShieldBytes           [][]uint64                      `json:"shield_bytes,omitempty"`
				CDNBytes              [][]uint64                      `json:"cdn_bytes,omitempty"`
				UpstreamBytes         [][]uint64                      `json:"upstream_bytes,omitempty"`
				Responses2xx          [][]uint64                      `json:"responses_2xx,omitempty"`
				Responses3xx          [][]uint64                      `json:"responses_3xx,omitempty"`
				Responses4xx          [][]uint64                      `json:"responses_4xx,omitempty"`
				Responses5xx          [][]uint64                      `json:"responses_5xx,omitempty"`
				ResponsesHit          [][]uint64                      `json:"responses_hit,omitempty"`
				ResponsesMiss         [][]uint64                      `json:"responses_miss,omitempty"`
				ImageProcessed        [][]uint64                      `json:"image_processed,omitempty"`
				OriginResponseTime    [][]uint64                      `json:"origin_response_time,omitempty"`
				RequestTime           [][]uint64                      `json:"request_time,omitempty"`
				CacheHitTrafficRatio  [][]float64                     `json:"cache_hit_traffic_ratio,omitempty"`
				CacheHitRequestsRatio [][]float64                     `json:"cache_hit_requests_ratio,omitempty"`
				ShieldTrafficRatio    [][]float64                     `json:"shield_traffic_ratio,omitempty"`
				RequestsWafPassed     [][]float64                     `json:"requests_waf_passed,omitempty"`
				ShieldUsage           []ResourceStatisticsShieldUsage `json:"shield_usage,omitempty"`
			}{Requests: [][]uint64{{1704067200, 500}}},
		},
	}

	var capturedPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.RequestURI()
		assert.Equal(t, http.MethodGet, r.Method)
		json.NewEncoder(w).Encode(expected)
	}))
	defer ts.Close()

	req := &ResourceStatisticsTimeSeriesRequest{
		Metrics:     []Metric{MetricRequests},
		GroupBy:     []GroupBy{GroupByRegion},
		Regions:     []Region{RegionEU},
		Granularity: Granularity1h,
		From:        from,
		To:          to,
	}

	service := NewService(provider.NewClient(ts.URL))
	result, err := service.GetTimeSeriesData(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, &expected, result)
	assert.Contains(t, capturedPath, "granularity=1h")
	assert.Contains(t, capturedPath, "group_by=region")
	assert.Contains(t, capturedPath, "metrics=requests")
	assert.Contains(t, capturedPath, "regions=eu")
}

func TestResourceStatisticsService_GroupByCombinations(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		groupBy        []GroupBy
		expectedParams []string
	}{
		{
			name:           "group by resource",
			groupBy:        []GroupBy{GroupByResource},
			expectedParams: []string{"group_by=resource"},
		},
		{
			name:           "group by region",
			groupBy:        []GroupBy{GroupByRegion},
			expectedParams: []string{"group_by=region"},
		},
		{
			name:           "group by country",
			groupBy:        []GroupBy{GroupByCountry},
			expectedParams: []string{"group_by=country"},
		},
		{
			name:           "multiple group by",
			groupBy:        []GroupBy{GroupByResource, GroupByRegion},
			expectedParams: []string{"group_by=resource", "group_by=region"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedPath string
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedPath = r.URL.RequestURI()
				json.NewEncoder(w).Encode(ResourceStatisticsTableResponse{})
			}))
			defer ts.Close()

			service := NewService(provider.NewClient(ts.URL))
			_, err := service.GetTableData(context.Background(), &ResourceStatisticsTableRequest{
				Metrics: []Metric{MetricRequests},
				GroupBy: tt.groupBy,
				From:    from,
				To:      to,
			})

			require.NoError(t, err)
			for _, param := range tt.expectedParams {
				assert.Contains(t, capturedPath, param)
			}
		})
	}
}

func TestResourceStatisticsService_GranularityCombinations(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		granularity Granularity
	}{
		{name: "1m", granularity: Granularity1m},
		{name: "5m", granularity: Granularity5m},
		{name: "15m", granularity: Granularity15m},
		{name: "1h", granularity: Granularity1h},
		{name: "1d", granularity: Granularity1d},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedPath string
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedPath = r.URL.RequestURI()
				json.NewEncoder(w).Encode(ResourceStatisticsTimeSeriesResponse{})
			}))
			defer ts.Close()

			service := NewService(provider.NewClient(ts.URL))
			_, err := service.GetTimeSeriesData(context.Background(), &ResourceStatisticsTimeSeriesRequest{
				Metrics:     []Metric{MetricRequests},
				GroupBy:     []GroupBy{GroupByResource},
				Granularity: tt.granularity,
				From:        from,
				To:          to,
			})

			require.NoError(t, err)
			assert.Contains(t, capturedPath, "granularity="+string(tt.granularity))
		})
	}
}

func TestResourceStatisticsService_GetTimeSeriesData_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	_, err := service.GetTimeSeriesData(context.Background(), &ResourceStatisticsTimeSeriesRequest{
		Metrics:     []Metric{MetricRequests},
		Granularity: Granularity1h,
		From:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		To:          time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestResourceStatisticsService_BadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"Message": "invalid parameters"})
	}))
	defer ts.Close()

	service := NewService(provider.NewClient(ts.URL))
	_, err := service.GetTableData(context.Background(), &ResourceStatisticsTableRequest{
		From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid parameters")
}
