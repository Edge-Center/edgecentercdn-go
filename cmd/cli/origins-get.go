package cli

import (
	"github.com/spf13/cobra"
	"os"
)

var getOriginCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an origin group",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetInt64("id")

		client, err := NewServiceCommandCobra(cmd)
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		result, err := client.OriginGroups().Get(ctx, id)

		if err != nil {
			return err
		}

		return PrintAsJSON(os.Stdout, result)
	},
}

func init() {
	getOriginCmd.Flags().Int64("id", 0, "id of the origin group")
	_ = getOriginCmd.MarkFlagRequired("id")

	originsCmd.AddCommand(getOriginCmd)
}
