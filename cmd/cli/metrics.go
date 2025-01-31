package cli

import (
	"github.com/Edge-Center/edgecentercdn-go/statistics"
	"github.com/spf13/cobra"
	"os"
)

type PrintFormat = string

const (
	PrintFormatJSON  PrintFormat = "json"
	PrintFormatTable PrintFormat = "table"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Work with traffic metrics",
}

var getMetricsCmd = &cobra.Command{
	Use:   "get",
	Short: "Get traffic metrics",
	RunE: func(cmd *cobra.Command, args []string) error {
		metrics, _ := cmd.Flags().GetStringSlice("metric")
		regions, _ := cmd.Flags().GetStringSlice("region")
		groupBy, _ := cmd.Flags().GetStringSlice("groupby")
		granularity, _ := cmd.Flags().GetString("granularity")
		hosts, _ := cmd.Flags().GetStringSlice("host")
		resources, _ := cmd.Flags().GetInt64Slice("resource")
		countries, _ := cmd.Flags().GetStringSlice("country")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		output, _ := cmd.Flags().GetString("output")

		if err := validateMetrics(metrics); err != nil {
			return err
		}
		if err := validateRegions(regions); err != nil {
			return err
		}
		if err := validateGroupBy(groupBy); err != nil {
			return err
		}
		if err := validateGranularity(granularity); err != nil {
			return err
		}
		if err := validateTimeRange(from, to); err != nil {
			return err
		}

		fromTime, toTime, err := statistics.ParseTimeRange(from, to)
		if err != nil {
			return err
		}

		client, err := NewServiceCommandCobra(cmd)
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		if granularity != "" {
			req := &statistics.ResourceStatisticsTimeSeriesRequest{
				Metrics:     convertMetrics(metrics),
				Regions:     convertRegions(regions),
				GroupBy:     convertGroupBy(groupBy),
				Granularity: convertGranularity(granularity),
				Hosts:       convertHosts(hosts),
				Resources:   convertResources(resources),
				Countries:   convertCountries(countries),
				From:        fromTime,
				To:          toTime,
			}

			result, err := client.Statistics().GetTimeSeriesData(ctx, req)
			if err != nil {
				return err
			}

			return printer(result)
		}

		req := &statistics.ResourceStatisticsTableRequest{
			Metrics:   convertMetrics(metrics),
			Regions:   convertRegions(regions),
			GroupBy:   convertGroupBy(groupBy),
			Hosts:     convertHosts(hosts),
			Resources: convertResources(resources),
			Countries: convertCountries(countries),
			From:      fromTime,
			To:        toTime,
		}

		result, err := client.Statistics().GetTableData(ctx, req)
		if err != nil {
			return err
		}

		if output == PrintFormatTable {
			return PrintStatisticsTable(os.Stdout, result, req.Metrics, req.GroupBy)
		}

		return printer(result)
	},
}

func init() {
	getMetricsCmd.Flags().String("from", "", "start time (e.g., 2025-01-01T00:00:00Z)")
	_ = getMetricsCmd.MarkFlagRequired("from")

	getMetricsCmd.Flags().String("to", "", "end time (e.g., 2025-01-15T23:59:00Z)")
	_ = getMetricsCmd.MarkFlagRequired("to")

	getMetricsCmd.Flags().StringSlice("metric", []string{}, "metrics to fetch")
	_ = getMetricsCmd.MarkFlagRequired("metric")

	getMetricsCmd.Flags().String("granularity", "", "granularity of data (1h, 1d)")
	getMetricsCmd.Flags().StringSlice("groupby", []string{}, "group by field (e.g., region, vhost, country, client, resource)")
	getMetricsCmd.Flags().StringSlice("region", []string{}, "groups by the regions")
	getMetricsCmd.Flags().StringSlice("host", []string{}, "groups by the custom domain of the CDN resource")
	getMetricsCmd.Flags().StringSlice("country", []string{}, "groups by the countries")
	getMetricsCmd.Flags().Int64Slice("resource", []int64{}, "groups by the CDN resource ID")
	getMetricsCmd.Flags().StringP("output", "o", PrintFormatJSON, "specify the output format (e.g., table, json)")

	_ = getMetricsCmd.RegisterFlagCompletionFunc("metric", cobra.FixedCompletions(statistics.MetricsSuggestions, cobra.ShellCompDirectiveNoFileComp))
	_ = getMetricsCmd.RegisterFlagCompletionFunc("groupby", cobra.FixedCompletions(statistics.GroupBySuggestions, cobra.ShellCompDirectiveNoFileComp))
	_ = getMetricsCmd.RegisterFlagCompletionFunc("granularity", cobra.FixedCompletions(statistics.GranularitySuggestions, cobra.ShellCompDirectiveNoFileComp))
	_ = getMetricsCmd.RegisterFlagCompletionFunc("region", cobra.FixedCompletions(statistics.RegionSuggestions, cobra.ShellCompDirectiveNoFileComp))

	metricsCmd.AddCommand(getMetricsCmd)
	rootCmd.AddCommand(metricsCmd)
}
