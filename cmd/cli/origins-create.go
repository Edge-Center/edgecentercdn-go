package cli

import (
	"encoding/json"
	"fmt"
	"github.com/Edge-Center/edgecentercdn-go/origingroups"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var createOriginCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new origin group",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := NewServiceCommandCobra(cmd)
		if err != nil {
			return err
		}

		input, _ := cmd.Flags().GetBool("input")

		if input {
			req := &origingroups.GroupRequest{}
			inputReader := cmd.InOrStdin()
			if err := json.NewDecoder(inputReader).Decode(req); err != nil {
				return err
			}

			ctx := cmd.Context()

			result, err := client.OriginGroups().Create(ctx, req)

			if err != nil {
				return err
			}

			return PrintAsJSON(os.Stdout, result)
		}

		name, _ := cmd.Flags().GetString("name")
		useNext, _ := cmd.Flags().GetBool("use-next")
		consistentBalancing, _ := cmd.Flags().GetBool("consistent-balancing")
		sources, _ := cmd.Flags().GetStringSlice("source")

		var origins []origingroups.OriginRequest

		for _, source := range sources {
			parts := strings.Split(source, ",")

			if len(parts) == 1 {
				origins = append(origins, origingroups.OriginRequest{
					Source:  removeScheme(parts[0]),
					Backup:  false,
					Enabled: true,
				})
			} else if len(parts) == 3 {
				if parts[1] != "true" && parts[1] != "false" {
					return fmt.Errorf("invalid backup value: %s, expected 'true' or 'false'", parts[1])
				}
				if parts[2] != "true" && parts[2] != "false" {
					return fmt.Errorf("invalid enabled value: %s, expected 'true' or 'false'", parts[2])
				}

				backup := parts[1] == "true"
				enabled := parts[2] == "true"
				origins = append(origins, origingroups.OriginRequest{
					Source:  removeScheme(parts[0]),
					Backup:  backup,
					Enabled: enabled,
				})
			} else {
				return fmt.Errorf("invalid source format: %s, expected format 'source' or 'source,backup,enabled'", source)
			}
		}

		authType, _ := cmd.Flags().GetString("auth-type")
		accessKeyID, _ := cmd.Flags().GetString("access-key-id")
		secretKey, _ := cmd.Flags().GetString("secret-key")
		bucketName, _ := cmd.Flags().GetString("bucket-name")

		var authorization *origingroups.Authorization
		if authType != "" || accessKeyID != "" || secretKey != "" || bucketName != "" {
			authorization = &origingroups.Authorization{
				AuthType:    authType,
				AccessKeyID: accessKeyID,
				SecretKey:   secretKey,
				BucketName:  bucketName,
			}
		}

		req := &origingroups.GroupRequest{
			Name:                name,
			UseNext:             useNext,
			Origins:             origins,
			Authorization:       authorization,
			ConsistentBalancing: consistentBalancing,
		}

		ctx := cmd.Context()

		result, err := client.OriginGroups().Create(ctx, req)

		if err != nil {
			return err
		}

		return PrintAsJSON(os.Stdout, result)
	},
}

func init() {
	createOriginCmd.Flags().BoolP("input", "i", false, "input file")
	createOriginCmd.Flags().String("name", "", "name of the origin group")
	createOriginCmd.Flags().Bool("use-next", true, "enable 'useNext' feature")
	createOriginCmd.Flags().Bool("consistent-balancing", true, "enable consistent balancing")
	createOriginCmd.Flags().StringSlice("source", []string{}, "source in the format 'source' or 'source,backup,enabled'")
	createOriginCmd.Flags().String("auth-type", "", "authorization type (e.g., aws_signature_v2, aws_signature_v4)")
	createOriginCmd.Flags().String("access-key-id", "", "access key ID")
	createOriginCmd.Flags().String("secret-key", "", "secret key")
	createOriginCmd.Flags().String("bucket-name", "", "bucket name")

	originsCmd.AddCommand(createOriginCmd)
}
