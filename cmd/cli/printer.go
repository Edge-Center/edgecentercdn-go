package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Edge-Center/edgecentercdn-go/statistics"
	"github.com/dustin/go-humanize"
	"io"
	"math"
	"text/tabwriter"
)

func PrintAsJSON(writer io.Writer, data interface{}) error {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	_, err = writer.Write(prettyJSON)
	return err
}

func PrintStatisticsTable(writer io.Writer, response *statistics.ResourceStatisticsTableResponse, metrics []statistics.Metric, groupBy []statistics.GroupBy) error {
	if response == nil {
		return errors.New("response is nil")
	}

	tw := tabwriter.NewWriter(writer, 0, 0, 2, ' ', tabwriter.AlignRight)

	writeFormatted := func(writer io.Writer, format string, args ...interface{}) error {
		_, err := fmt.Fprintf(writer, format, args...)
		return err
	}

	for _, group := range groupBy {
		if err := writeFormatted(tw, "%s\t", statistics.GetGroupByLabel(group)); err != nil {
			return err
		}
	}

	for _, metric := range metrics {
		if err := writeFormatted(tw, "%s\t", statistics.GetMetricLabel(metric)); err != nil {
			return err
		}
	}

	if err := writeFormatted(tw, "\n"); err != nil {
		return err
	}

	for _, data := range *response {
		for _, group := range groupBy {
			var groupValue string
			switch group {
			case statistics.GroupByResource:
				if data.Resource != nil {
					groupValue = fmt.Sprintf("%.0f", *data.Resource)
				}
			case statistics.GroupByRegion:
				if data.Region != nil {
					groupValue = *data.Region
				}
			case statistics.GroupByHost:
				if data.Host != nil {
					groupValue = *data.Host
				}
			case statistics.GroupByClient:
				if data.Client != nil {
					groupValue = fmt.Sprintf("%.0f", *data.Client)
				}
			case statistics.GroupByClientRegion:
				if data.ClientRegion != nil {
					groupValue = *data.ClientRegion
				}
			case statistics.GroupByCountry:
				if data.Country != nil {
					groupValue = *data.Country
				}
			case statistics.GroupByDatacenter:
				if data.Datacenter != nil {
					groupValue = *data.Datacenter
				}
			}

			if err := writeFormatted(tw, "%s\t", groupValue); err != nil {
				return err
			}
		}

		for _, metric := range metrics {
			var value string
			switch metric {
			case statistics.MetricRequests:
				value = formatUint64WithCommas(data.Metrics.Requests)
			case statistics.MetricSentBytes:
				value = formatBytes(data.Metrics.SentBytes)
			case statistics.MetricShieldBytes:
				value = formatBytes(data.Metrics.ShieldBytes)
			case statistics.MetricTotalBytes:
				value = formatBytes(data.Metrics.TotalBytes)
			case statistics.MetricUpstreamBytes:
				value = formatBytes(data.Metrics.UpstreamBytes)
			case statistics.MetricCdnBytes:
				value = formatBytes(data.Metrics.CDNBytes)
			case statistics.MetricResponses2xx:
				value = formatUint64WithCommas(data.Metrics.Responses2xx)
			case statistics.MetricResponses3xx:
				value = formatUint64WithCommas(data.Metrics.Responses3xx)
			case statistics.MetricResponses4xx:
				value = formatUint64WithCommas(data.Metrics.Responses4xx)
			case statistics.MetricResponses5xx:
				value = formatUint64WithCommas(data.Metrics.Responses5xx)
			case statistics.MetricResponsesHit:
				value = formatUint64WithCommas(data.Metrics.ResponsesHit)
			case statistics.MetricResponsesMiss:
				value = formatUint64WithCommas(data.Metrics.ResponsesMiss)
			case statistics.MetricImageProcessed:
				value = formatUint64WithCommas(data.Metrics.ImageProcessed)
			case statistics.MetricRequestWafPassed:
				value = formatUint64WithCommas(data.Metrics.RequestsWafPassed)
			case statistics.MetricCacheHitTrafficRatio:
				value = formatPercentage(data.Metrics.CacheHitTrafficRatio)
			case statistics.MetricCacheHitRequestsRatio:
				value = formatPercentage(data.Metrics.CacheHitRequestsRatio)
			case statistics.MetricShieldTrafficRatio:
				value = formatPercentage(data.Metrics.ShieldTrafficRatio)
			case statistics.MetricOriginResponseTime:
				value = formatSeconds(data.Metrics.OriginResponseTime)
			case statistics.MetricRequestTime:
				value = formatSeconds(data.Metrics.RequestTime)
			}
			if err := writeFormatted(tw, "%s\t", value); err != nil {
				return err
			}
		}
		if err := writeFormatted(tw, "\n"); err != nil {
			return err
		}
	}

	return tw.Flush()
}

func formatUint64WithCommas(num uint64) string {
	return humanize.Commaf(float64(num))
}

func formatBytes(num uint64) string {
	return humanize.Bytes(num)
}

func formatPercentage(value float64) string {
	percentage := math.Round(value*1000) / 10
	return fmt.Sprintf("%.1f%%", percentage)
}

func formatSeconds(nanos uint64) string {
	switch {
	case nanos < 1_000:
		return fmt.Sprintf("%d ns", nanos)
	case nanos < 1_000_000:
		return fmt.Sprintf("%.2f µs", float64(nanos)/1_000)
	case nanos < 1_000_000_000:
		return fmt.Sprintf("%.2f ms", float64(nanos)/1_000_000)
	default:
		seconds := float64(nanos) / 1_000_000_000
		return fmt.Sprintf("%.2f s", seconds)
	}
}
