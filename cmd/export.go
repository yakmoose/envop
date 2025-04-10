/*
Copyright © 2025 John Lennard <john@yakmoo.se>
*/
package cmd

import (
	"context"
	"github.com/1password/onepassword-sdk-go"
	"github.com/spf13/cobra"
	"github.com/yakmoose/envop/service"
	"strings"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export the specified 1password item into an environment file",
	RunE: func(cmd *cobra.Command, args []string) error {
		envFile, err := cmd.Flags().GetString("env-file")
		if err != nil {
			return err
		}

		vaultName, err := cmd.Flags().GetString("vault")
		if err != nil {
			return err
		}

		itemName, err := cmd.Flags().GetString("item")
		if err != nil {
			return err
		}

		token, err := cmd.Flags().GetString("service-account")
		if err != nil {
			return err
		}

		client, err := onepassword.NewClient(
			context.Background(),
			onepassword.WithServiceAccountToken(token),
			onepassword.WithIntegrationInfo("envop", "v0.0.0"),
		)
		if err != nil {
			return err
		}

		item, err := service.Get1PasswordItem(client, vaultName, itemName)
		if err != nil {
			return err
		}

		fields := make(map[string]string, 0)
		for _, v := range item.Fields {
			fields[strings.ToUpper(v.Title)] = v.Value
		}

		format, err := cmd.Flags().GetString("format")
		if err != nil {
			return err
		}

		switch format {
		case "json":
			service.WriteJSON(envFile, fields)
		case "env":
			service.WriteEnv(envFile, fields)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().String("env-file", "", "The file to save to")

	exportCmd.Flags().String("vault", "", "The 1password vault")
	exportCmd.MarkFlagRequired("vault")

	exportCmd.Flags().String("item", "", "The name of the item to save")
	exportCmd.MarkFlagRequired("item")

	exportCmd.Flags().String("format", "env", "The file format to save as (env, json)")

}
