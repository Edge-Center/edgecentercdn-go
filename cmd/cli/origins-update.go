package cli

import (
	"encoding/json"
	"fmt"
	"github.com/Edge-Center/edgecentercdn-go/origingroups"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var updateOriginCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an origin group",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := NewServiceCommandCobra(cmd)
		if err != nil {
			return err
		}

		id, _ := cmd.Flags().GetInt64("id")
		input, _ := cmd.Flags().GetBool("input")

		if input {
			req := &origingroups.GroupRequest{}
			inputReader := cmd.InOrStdin()
			if err := json.NewDecoder(inputReader).Decode(req); err != nil {
				return err
			}

			ctx := cmd.Context()

			result, err := client.OriginGroups().Update(ctx, id, req)

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

		result, err := client.OriginGroups().Update(ctx, id, req)

		if err != nil {
			return err
		}

		return PrintAsJSON(os.Stdout, result)
	},
}

func init() {
	updateOriginCmd.Flags().BoolP("input", "i", false, "input file")
	updateOriginCmd.Flags().Int64("id", 0, "id of the origin group")
	updateOriginCmd.Flags().String("name", "", "name of the origin group")
	updateOriginCmd.Flags().Bool("use-next", true, "enable 'useNext' feature")
	updateOriginCmd.Flags().Bool("consistent-balancing", true, "enable consistent balancing")
	updateOriginCmd.Flags().StringSlice("source", []string{}, "source in the format 'source' or 'source,backup,enabled'")
	updateOriginCmd.Flags().String("auth-type", "", "authorization type (e.g., aws_signature_v2, aws_signature_v4)")
	updateOriginCmd.Flags().String("access-key-id", "", "access key ID")
	updateOriginCmd.Flags().String("secret-key", "", "secret key")
	updateOriginCmd.Flags().String("bucket-name", "", "bucket name")

	_ = updateOriginCmd.MarkFlagRequired("name")
	_ = updateOriginCmd.MarkFlagRequired("id")

	originsCmd.AddCommand(updateOriginCmd)
}
