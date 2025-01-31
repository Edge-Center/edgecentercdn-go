package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var AppVersion = "v0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(AppVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
