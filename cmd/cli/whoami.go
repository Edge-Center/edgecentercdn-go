package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Displays information about the current user",
	RunE: func(cmd *cobra.Command, args []string) error {

		client, err := NewServiceCommandCobra(cmd)
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		username, err := client.Tools().Whoami(ctx)
		if err != nil {
			return err
		}

		fmt.Println(username)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
