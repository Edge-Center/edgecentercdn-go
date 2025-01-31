package statistics

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func buildStatisticsPath(
	from time.Time,
	to time.Time,
	granularity Granularity,
	groupBy []GroupBy,
	metrics []Metric,
	regions []Region,
	hosts []Host,
	resources []ResourceId,
	clients []ClientId,
	countries []Country,
) string {
	baseURL := "/cdn/statistics/aggregate/stats"

	params := url.Values{}

	params.Set("service", "CDN")
	params.Set("flat", "true")
	params.Set("from", from.Format(time.RFC3339))
	params.Set("to", to.Format(time.RFC3339))

	if granularity != "" {
		params.Set("granularity", string(granularity))
	}

	for _, gr := range groupBy {
		params.Add("group_by", string(gr))
	}

	for _, metric := range metrics {
		params.Add("metrics", string(metric))
	}

	for _, region := range regions {
		params.Add("regions", string(region))
	}

	for _, clientID := range clients {
		params.Add("client", strconv.FormatInt(int64(clientID), 10))
	}

	for _, host := range hosts {
		params.Add("vhost", string(host))
	}

	for _, country := range countries {
		params.Add("countries", string(country))
	}

	for _, resourceID := range resources {
		params.Add("resource", strconv.FormatInt(int64(resourceID), 10))
	}

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

func ParseTimeRange(from, to string) (time.Time, time.Time, error) {
	var fromTime, toTime time.Time
	var err error
	if from != "" {
		fromTime, err = time.Parse(time.RFC3339, from)
		if err != nil {
			return fromTime, toTime, fmt.Errorf("invalid 'from' time format: %v", err)
		}
	}
	if to != "" {
		toTime, err = time.Parse(time.RFC3339, to)
		if err != nil {
			return fromTime, toTime, fmt.Errorf("invalid 'to' time format: %v", err)
		}
	}
	return fromTime, toTime, nil
}

func GranularityToSeconds(g Granularity) (int64, error) {
	switch g {
	case Granularity1m:
		return 60, nil
	case Granularity5m:
		return 300, nil
	case Granularity15m:
		return 900, nil
	case Granularity1h:
		return 3600, nil
	case Granularity1d:
		return 86400, nil
	default:
		return 0, errors.New("unknown granularity")
	}
}
