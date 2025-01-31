package cli

import (
	"github.com/Edge-Center/edgecentercdn-go/resources"
	"github.com/spf13/cobra"
	"os"
)

var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "Manage resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var listResourceCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		ordering, _ := cmd.Flags().GetString("ordering")
		search, _ := cmd.Flags().GetString("search")
		deleted, _ := cmd.Flags().GetBool("deleted")
		fields, _ := cmd.Flags().GetStringArray("field")
		statuses, _ := cmd.Flags().GetStringArray("status")

		req := &resources.ListFilterRequest{
			Ordering: ordering,
			Search:   search,
			Fields:   fields,
			Deleted:  deleted,
			Status:   convertStatus(statuses),
		}

		client, err := NewServiceCommandCobra(cmd)
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		result, err := client.Resources().List(ctx, req)

		if err != nil {
			return err
		}

		return PrintAsJSON(os.Stdout, result)
	},
}

func init() {
	listResourceCmd.Flags().String("ordering", "", "order the results")
	listResourceCmd.Flags().String("search", "", "search term")
	listResourceCmd.Flags().Bool("deleted", false, "include deleted resources in the list")
	listResourceCmd.Flags().StringArray("field", []string{}, "fields to include, comma-separated")
	listResourceCmd.Flags().StringArray("status", []string{}, "resource status (active, suspended, processed)")

	resourceCmd.AddCommand(listResourceCmd)

	rootCmd.AddCommand(resourceCmd)
}

func convertStatus(status []string) []resources.ResourceStatus {
	result := make([]resources.ResourceStatus, len(status))
	for i, m := range status {
		result[i] = resources.ResourceStatus(m)
	}
	return result
}
