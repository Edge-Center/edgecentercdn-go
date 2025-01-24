package cli

import (
	"fmt"
	"github.com/Edge-Center/edgecentercdn-go/resources"
	"github.com/Edge-Center/edgecentercdn-go/tools"
	"github.com/spf13/cobra"
)

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge specified paths for a given resource ID or cname",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetInt64("id")
		cname, _ := cmd.Flags().GetString("cname")
		paths, _ := cmd.Flags().GetStringSlice("path")

		client, err := NewServiceCommandCobra(cmd)
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		req := &tools.PurgeRequest{Paths: paths}

		if id != 0 {
			resp, err := client.Tools().Purge(ctx, id, req)
			if err != nil {
				return fmt.Errorf("failed to purge: %w", err)
			}

			if len(resp.Paths) == 0 {
				fmt.Printf("Successfully purged all files for resource ID %d\n", id)
			} else {
				fmt.Printf("Successfully purged paths for resource ID %d\n", id)
			}
		}

		if cname != "" {
			listReq := &resources.ListFilterRequest{
				Cname:  cname,
				Status: []resources.ResourceStatus{resources.ActiveResourceStatus},
				Fields: []string{"cname", "id"},
			}

			cdnResources, err := client.Resources().List(ctx, listReq)
			if err != nil {
				return fmt.Errorf("failed to list resources for cname %s: %w", cname, err)
			}

			if len(cdnResources) == 0 {
				return fmt.Errorf("no resources found for cname %s", cname)
			}

			for _, cdnResource := range cdnResources {
				resp, err := client.Tools().Purge(ctx, cdnResource.ID, req)
				if err != nil {
					fmt.Printf("Error purging for resource ID %d (cname: %s): %v\n", cdnResource.ID, cdnResource.Cname, err)
					continue
				}

				if len(resp.Paths) == 0 {
					fmt.Printf("Successfully purged all files for resource ID %d (cname: %s)\n", cdnResource.ID, cdnResource.Cname)
				} else {
					fmt.Printf("Successfully purged paths for resource ID %d (cname: %s)\n", cdnResource.ID, cdnResource.Cname)
				}
			}
		}

		return nil
	},
}

func init() {
	purgeCmd.Flags().Int64("id", 0, "Resource ID")
	purgeCmd.Flags().String("cname", "", "Resource cname")
	purgeCmd.Flags().StringSliceP("path", "p", []string{}, "Paths to purge")

	rootCmd.AddCommand(purgeCmd)
}
