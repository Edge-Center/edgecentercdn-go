package cli

import (
	"fmt"
	"github.com/Edge-Center/edgecentercdn-go/statistics"
	"github.com/spf13/cobra"
)

var manCmd = &cobra.Command{
	Use:   "man",
	Short: "Display information about metrics and groupings",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var manMetricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Display information about available metrics",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available Metrics:")
		fmt.Println("------------------")
		for metric, details := range statistics.Strings.Metrics {
			fmt.Printf("%s\n  Label: %s\n  Description: %s\n\n", metric, details.Label, details.Desc)
		}
	},
}

var manGroupByCmd = &cobra.Command{
	Use:   "groupBy",
	Short: "Display information about available groupings",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available GroupBy Options:")
		fmt.Println("--------------------------")
		for groupBy, details := range statistics.Strings.GroupBy {
			fmt.Printf("%s\n  Label: %s\n  Description: %s\n\n", groupBy, details.Label, details.Desc)
		}
	},
}

var manGranularityByCmd = &cobra.Command{
	Use:   "granularity",
	Short: "Display information about available granularity",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available granularity Options:")
		fmt.Println("--------------------------")
		for granularity, details := range statistics.Strings.Granularity {
			fmt.Printf("%s\n  Label: %s\n  Description: %s\n\n", granularity, details.Label, details.Desc)
		}
	},
}

func init() {
	manCmd.AddCommand(manMetricsCmd, manGroupByCmd, manGranularityByCmd)

	rootCmd.AddCommand(manCmd)
}
