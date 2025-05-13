package statistics

import (
	"context"
	"time"
)

type ResourceStatisticsService interface {
	GetTableData(ctx context.Context, req *ResourceStatisticsTableRequest) (*ResourceStatisticsTableResponse, error)
	GetTimeSeriesData(ctx context.Context, req *ResourceStatisticsTimeSeriesRequest) (*ResourceStatisticsTimeSeriesResponse, error)
}

type Region string
type Metric string
type GroupBy string
type Granularity string
type Host string
type Country string
type ResourceId int64
type ClientId int64

const (
	MetricUpstreamBytes         Metric = "upstream_bytes"
	MetricSentBytes             Metric = "sent_bytes"
	MetricShieldBytes           Metric = "shield_bytes"
	MetricTotalBytes            Metric = "total_bytes"
	MetricCdnBytes              Metric = "cdn_bytes"
	MetricRequestTime           Metric = "request_time"
	MetricRequests              Metric = "requests"
	MetricResponses2xx          Metric = "responses_2xx"
	MetricResponses3xx          Metric = "responses_3xx"
	MetricResponses4xx          Metric = "responses_4xx"
	MetricResponses5xx          Metric = "responses_5xx"
	MetricResponsesHit          Metric = "responses_hit"
	MetricResponsesMiss         Metric = "responses_miss"
	MetricRequestWafPassed      Metric = "requests_waf_passed"
	MetricCacheHitTrafficRatio  Metric = "cache_hit_traffic_ratio"
	MetricCacheHitRequestsRatio Metric = "cache_hit_requests_ratio"
	MetricShieldTrafficRatio    Metric = "shield_traffic_ratio"
	MetricImageProcessed        Metric = "image_processed"
	MetricOriginResponseTime    Metric = "origin_response_time"
)

const (
	GroupByResource     GroupBy = "resource"
	GroupByRegion       GroupBy = "region"
	GroupByHost         GroupBy = "vhost"
	GroupByClient       GroupBy = "client"
	GroupByCountry      GroupBy = "country"
	GroupByClientRegion GroupBy = "client_region"
	GroupByDatacenter   GroupBy = "dc"
)

const (
	RegionNA     Region = "na"
	RegionEU     Region = "eu"
	RegionCIS    Region = "cis"
	RegionAsia   Region = "asia"
	RegionAU     Region = "au"
	RegionLATAM  Region = "latam"
	RegionME     Region = "me"
	RegionAfrica Region = "africa"
	RegionSA     Region = "sa"
	RegionRU     Region = "ru"
)

const (
	Granularity1m  Granularity = "1m"
	Granularity5m  Granularity = "5m"
	Granularity15m Granularity = "15m"
	Granularity1h  Granularity = "1h"
	Granularity1d  Granularity = "1d"
)

type ResourceStatisticsTimeSeriesRequest struct {
	Metrics     []Metric
	Regions     []Region
	GroupBy     []GroupBy
	Granularity Granularity
	Hosts       []Host
	Resources   []ResourceId
	Countries   []Country
	Clients     []ClientId
	From        time.Time
	To          time.Time
}

type ResourceStatisticsTableRequest struct {
	Metrics   []Metric
	Regions   []Region
	GroupBy   []GroupBy
	Hosts     []Host
	Resources []ResourceId
	Clients   []ClientId
	Countries []Country
	From      time.Time
	To        time.Time
}

type ResourceStatisticsShieldUsage struct {
	ActiveFrom time.Time `json:"active_from"`
	ActiveTo   time.Time `json:"active_to"`
	ClientID   int       `json:"client_id"`
	Cname      string    `json:"cname"`
}

type ResourceStatisticsTimeSeriesData struct {
	Metrics struct {
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
	} `json:"metrics"`
	Resource     *float64 `json:"resource,omitempty"`
	Client       *float64 `json:"client,omitempty"`
	Region       *string  `json:"region,omitempty"`
	Country      *string  `json:"country,omitempty"`
	Datacenter   *string  `json:"dc,omitempty"`
	Host         *string  `json:"vhost,omitempty"`
	ClientRegion *string  `json:"client_region,omitempty"`
}

type ResourceStatisticsTimeSeriesResponse = []ResourceStatisticsTimeSeriesData

type ResourceStatisticsTableData struct {
	Metrics struct {
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
	} `json:"metrics"`
	Resource     *float64 `json:"resource,omitempty"`
	Client       *float64 `json:"client,omitempty"`
	Region       *string  `json:"region,omitempty"`
	Country      *string  `json:"country,omitempty"`
	Datacenter   *string  `json:"dc,omitempty"`
	Host         *string  `json:"vhost,omitempty"`
	ClientRegion *string  `json:"client_region,omitempty"`
}

type ResourceStatisticsTableResponse = []ResourceStatisticsTableData
