package cli

import (
	"fmt"
	"github.com/Edge-Center/edgecentercdn-go/statistics"
	"os"
	"strings"
)

func printer(data interface{}) error {
	return PrintAsJSON(os.Stdout, data)
}

func validateMetrics(metrics []string) error {
	if len(metrics) == 0 {
		return fmt.Errorf("at least one metric must be specified using --metric")
	}
	return statistics.ValidateMetrics(metrics)
}

func validateRegions(regions []string) error {
	if len(regions) == 0 {
		return nil
	}

	return statistics.ValidateRegions(regions)
}

func validateGroupBy(groupBy []string) error {
	if len(groupBy) == 0 {
		return nil
	}

	return statistics.ValidateGroupBy(groupBy)
}

func validateGranularity(granularity string) error {
	if len(granularity) == 0 {
		return nil
	}

	return statistics.ValidateGranularity(granularity)
}

func validateTimeRange(from, to string) error {
	return statistics.ValidateTimeRange(from, to)
}

func convertMetrics(metrics []string) []statistics.Metric {
	result := make([]statistics.Metric, len(metrics))
	for i, m := range metrics {
		result[i] = statistics.Metric(m)
	}
	return result
}

func convertRegions(regions []string) []statistics.Region {
	result := make([]statistics.Region, len(regions))
	for i, r := range regions {
		result[i] = statistics.Region(r)
	}
	return result
}

func convertCountries(countries []string) []statistics.Country {
	result := make([]statistics.Country, len(countries))
	for i, r := range countries {
		result[i] = statistics.Country(r)
	}
	return result
}

func convertGroupBy(groupBy []string) []statistics.GroupBy {
	result := make([]statistics.GroupBy, len(groupBy))
	for i, r := range groupBy {
		result[i] = statistics.GroupBy(r)
	}
	return result
}

func convertGranularity(granularity string) statistics.Granularity {
	return statistics.Granularity(granularity)
}

func convertHosts(vhosts []string) []statistics.Host {
	result := make([]statistics.Host, len(vhosts))
	for i, r := range vhosts {
		result[i] = statistics.Host(r)
	}
	return result
}

func convertResources(resources []int64) []statistics.ResourceId {
	result := make([]statistics.ResourceId, len(resources))
	for i, r := range resources {
		result[i] = statistics.ResourceId(r)
	}
	return result
}

func removeScheme(url string) string {
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	return url
}
