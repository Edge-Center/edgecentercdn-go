package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var deleteOriginCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an origin group",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetInt64("id")

		client, err := NewServiceCommandCobra(cmd)
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		err = client.OriginGroups().Delete(ctx, id)

		if err != nil {
			return err
		}

		fmt.Printf("Successfully deleted origin group with ID: %d\n", id)
		return nil
	},
}

func init() {
	deleteOriginCmd.Flags().Int64("id", 0, "id of the origin group")
	_ = deleteOriginCmd.MarkFlagRequired("id")

	originsCmd.AddCommand(deleteOriginCmd)
}
