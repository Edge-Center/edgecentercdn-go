package statistics

type StringDetails struct {
	Label string `json:"label"`
	Desc  string `json:"desc"`
}

type StringMappings struct {
	Metrics     map[Metric]StringDetails      `json:"metrics"`
	GroupBy     map[GroupBy]StringDetails     `json:"groupBy"`
	Granularity map[Granularity]StringDetails `json:"granularity"`
}

var metricsMapping = map[Metric]StringDetails{
	MetricTotalBytes: {
		Label: "Total Traffic",
		Desc:  "The total traffic volume. It is the sum of the traffic from the origin to the CDN servers or shielding, the traffic from the shielding to the CDN servers, and the traffic from the CDN servers to the end users.",
	},
	MetricUpstreamBytes: {
		Label: "Origin Traffic",
		Desc:  "The traffic from the origin. It is the traffic from the origin to the CDN servers or from the origin to the shielding.",
	},
	MetricSentBytes: {
		Label: "Edges Traffic",
		Desc:  "The traffic from the CDN servers. It is the sum of the traffic from the shielding to the CDN servers and the traffic from the CDN servers to the end users.",
	},
	MetricShieldBytes: {
		Label: "Shield Traffic",
		Desc:  "The traffic from the shielding. It is the traffic from the shielding to the CDN servers.",
	},
	MetricRequests: {
		Label: "Total Requests",
		Desc:  "The total number of requests from the end users to the CDN servers.",
	},
	MetricRequestTime: {
		Label: "Request Time",
		Desc:  "The time taken to process requests.",
	},
	MetricResponses2xx: {
		Label: "2xx Responses",
		Desc:  "The number of 2xx responses from the CDN servers.",
	},
	MetricResponses3xx: {
		Label: "3xx Responses",
		Desc:  "The number of 3xx responses from the CDN servers.",
	},
	MetricResponses4xx: {
		Label: "4xx Responses",
		Desc:  "The number of 4xx responses from the CDN servers.",
	},
	MetricResponses5xx: {
		Label: "5xx Responses",
		Desc:  "The number of 5xx responses from the CDN servers.",
	},
	MetricRequestWafPassed: {
		Label: "WAF requests",
		Desc:  "The number of requests that were processed by Basic WAF.",
	},
	MetricResponsesHit: {
		Label: "Cache Hit Ratio",
		Desc:  "The amount of cached content that is sent. It is calculated by the formula: the traffic of the cached content from the CDN servers (with the Cache: HIT HTTP header) divided by the number of requests to the CDN servers.",
	},
	MetricResponsesMiss: {
		Label: "Cache Miss",
		Desc:  "The number of requests that resulted in a cache miss.",
	},
	MetricCacheHitTrafficRatio: {
		Label: "Byte Cache Hit Ratio",
		Desc:  "The amount of cached traffic. It is calculated by the formula: one minus the traffic from the origin to the CDN servers or the shielding divided by the traffic from the CDN servers to the end users.",
	},
	MetricCacheHitRequestsRatio: {
		Label: "Cache Hit Ratio",
		Desc:  "The amount of cached content that is sent. It is calculated by the formula: the traffic of the cached content from the CDN servers (with the Cache: HIT HTTP header) divided by the number of requests to the CDN servers.",
	},
	MetricShieldTrafficRatio: {
		Label: "Shield Traffic Ratio",
		Desc:  "The efficiency of the shielding: how much more traffic is sent from the shielding rather than from the origin. It is calculated by the formula: (the traffic from the shielding to the CDN servers minus the traffic from the origin to the shielding) divided by the traffic from the shielding to the CDN servers.",
	},
	MetricImageProcessed: {
		Label: "Image Optimization",
		Desc:  "The number of images processed by the Image Stack.",
	},
	MetricOriginResponseTime: {
		Label: "Origin Response Time",
		Desc:  "The time taken for the origin server to respond.",
	},
	MetricCdnBytes: {
		Label: "CDN Traffic",
		Desc:  "The sum of sent_bytes and shield_bytes traffic.",
	},
}

var groupByMapping = map[GroupBy]StringDetails{
	GroupByResource: {
		Label: "Resource ID",
		Desc:  "Groups by the CDN resource ID.",
	},
	GroupByRegion: {
		Label: "Region",
		Desc:  "Groups by the regions.",
	},
	GroupByHost: {
		Label: "Host",
		Desc:  "Groups by the custom domain of the CDN resource.",
	},
	GroupByClient: {
		Label: "Client ID",
		Desc:  "Groups by the client’s ID.",
	},
	GroupByCountry: {
		Label: "Country",
		Desc:  "Groups by the countries.",
	},
	GroupByClientRegion: {
		Label: "Client Region",
		Desc:  "Groups by the client region.",
	},
	GroupByDatacenter: {
		Label: "Datacenter",
		Desc:  "Groups by the data centers.",
	},
}

var granularityMapping = map[Granularity]StringDetails{
	Granularity1m: {
		Label: "1 minute",
		Desc:  "Data points aggregated every minute",
	},
	Granularity5m: {
		Label: "5 minutes",
		Desc:  "Data points aggregated every 5 minutes",
	},
	Granularity15m: {
		Label: "15 minutes",
		Desc:  "Data points aggregated every 15 minutes",
	},
	Granularity1h: {
		Label: "1 hour",
		Desc:  "Data points aggregated every hour",
	},
	Granularity1d: {
		Label: "1 day",
		Desc:  "Data points aggregated daily",
	},
}

var Strings = StringMappings{
	Metrics:     metricsMapping,
	GroupBy:     groupByMapping,
	Granularity: granularityMapping,
}

func GetMetricLabel(metric Metric) string {
	info, exists := Strings.Metrics[metric]
	if !exists {
		return string(metric)
	}
	return info.Label
}

func GetGroupByLabel(gr GroupBy) string {
	info, exists := Strings.GroupBy[gr]
	if !exists {
		return string(gr)
	}
	return info.Label
}
