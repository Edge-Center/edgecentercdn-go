package statistics

import (
	"errors"
	"fmt"
	"time"
)

func ValidateMetrics(metrics []string) error {
	validMetrics := map[Metric]bool{
		MetricUpstreamBytes:         true,
		MetricSentBytes:             true,
		MetricShieldBytes:           true,
		MetricTotalBytes:            true,
		MetricCdnBytes:              true,
		MetricRequestTime:           true,
		MetricRequests:              true,
		MetricResponses2xx:          true,
		MetricResponses3xx:          true,
		MetricResponses4xx:          true,
		MetricResponses5xx:          true,
		MetricResponsesHit:          true,
		MetricResponsesMiss:         true,
		MetricCacheHitTrafficRatio:  true,
		MetricCacheHitRequestsRatio: true,
		MetricShieldTrafficRatio:    true,
		MetricImageProcessed:        true,
		MetricOriginResponseTime:    true,
	}
	for _, metric := range metrics {
		m := Metric(metric)
		if !validMetrics[m] {
			return fmt.Errorf("invalid metric: %s", metric)
		}
	}
	return nil
}

func ValidateGranularity(granularity string) error {
	validGranularities := map[Granularity]bool{
		Granularity1m:  true,
		Granularity5m:  true,
		Granularity15m: true,
		Granularity1h:  true,
		Granularity1d:  true,
	}

	g := Granularity(granularity)
	if !validGranularities[g] {
		return fmt.Errorf("invalid granularity: %s", granularity)
	}
	return nil
}

func ValidateGroupBy(groupBy []string) error {
	validGroups := map[GroupBy]bool{
		GroupByRegion:       true,
		GroupByResource:     true,
		GroupByHost:         true,
		GroupByCountry:      true,
		GroupByClient:       true,
		GroupByClientRegion: true,
		GroupByDatacenter:   true,
	}
	for _, target := range groupBy {
		g := GroupBy(target)
		if !validGroups[g] {
			return fmt.Errorf("invalid groupBy value: %s", groupBy)
		}
	}
	return nil
}

func ValidateRegions(regions []string) error {
	validRegions := map[Region]bool{
		RegionNA:     true,
		RegionEU:     true,
		RegionCIS:    true,
		RegionAsia:   true,
		RegionAU:     true,
		RegionLATAM:  true,
		RegionME:     true,
		RegionAfrica: true,
		RegionSA:     true,
		RegionRU:     true,
	}

	for _, region := range regions {
		r := Region(region)
		if !validRegions[r] {
			return fmt.Errorf("invalid region: %s", region)
		}
	}
	return nil
}

func ValidateTimeRange(from, to string) error {
	if from == "" || to == "" {
		return errors.New("'from' and 'to' must both be set")
	}

	if _, err := time.Parse(time.RFC3339, from); err != nil {
		return fmt.Errorf("invalid 'from' time format: %v", err)
	}

	if _, err := time.Parse(time.RFC3339, to); err != nil {
		return fmt.Errorf("invalid 'to' time format: %v", err)
	}

	return nil
}
