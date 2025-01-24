package statistics

var MetricsSuggestions = []string{
	string(MetricUpstreamBytes),
	string(MetricSentBytes),
	string(MetricShieldBytes),
	string(MetricTotalBytes),
	string(MetricCdnBytes),
	string(MetricRequestTime),
	string(MetricRequests),
	string(MetricResponses2xx),
	string(MetricResponses3xx),
	string(MetricResponses4xx),
	string(MetricResponses5xx),
	string(MetricResponsesHit),
	string(MetricResponsesMiss),
	string(MetricCacheHitTrafficRatio),
	string(MetricCacheHitRequestsRatio),
	string(MetricShieldTrafficRatio),
	string(MetricImageProcessed),
	string(MetricOriginResponseTime),
	string(MetricRequestWafPassed),
}

var GranularitySuggestions = []string{
	string(Granularity1m),
	string(Granularity5m),
	string(Granularity15m),
	string(Granularity1h),
	string(Granularity1d),
}

var RegionSuggestions = []string{
	string(RegionNA),
	string(RegionEU),
	string(RegionCIS),
	string(RegionAsia),
	string(RegionAU),
	string(RegionLATAM),
	string(RegionME),
	string(RegionAfrica),
	string(RegionSA),
	string(RegionRU),
}

var GroupBySuggestions = []string{
	string(GroupByResource),
	string(GroupByHost),
	string(GroupByRegion),
	string(GroupByClient),
	string(GroupByCountry),
	string(GroupByClientRegion),
}
