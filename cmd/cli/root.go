package cli

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "edge-cli",
	Short: "EdgeCenter CDN CLI client",
	Long: `EdgeCenter CDN CLI client

Please make sure you set the environment variables:

EC_CDN_API_KEY - The API key for CDN access
EC_CDN_API_URL - The URL for the CDN API

Example:

export EC_CDN_API_KEY=your-api-key
export EC_CDN_API_URL=https://api.edgecenter.ru`,
	SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if os.Geteuid() == 0 {
			panic("Running as root is not recommended. Please use a non-root user.")
		}
	},
}
