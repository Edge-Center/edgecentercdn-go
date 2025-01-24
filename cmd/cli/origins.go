package cli

import (
	"github.com/spf13/cobra"
)

var originsCmd = &cobra.Command{
	Use:   "origins",
	Short: "Manage origins",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(originsCmd)
}
