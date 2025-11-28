/*
Copyright © 2025 John Lennard <john@yakmoo.se>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yakmoose/envop/service"
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

		sectionName, err := cmd.Flags().GetString("section")
		if err != nil {
			return err
		}

		token, err := cmd.Flags().GetString("service-account")
		if err != nil {
			return err
		}

		client, err := service.NewClientFromToken(token)
		if err != nil {
			return err
		}

		env, err := service.ReadOnePassword(
			client,
			vaultName,
			itemName,
			sectionName,
		)

		if err != nil {
			return err
		}

		format, err := cmd.Flags().GetString("format")
		if err != nil {
			return err
		}

		switch format {
		case "json":
			err = service.WriteJSON(envFile, env)
		case "env":
			err = service.WriteEnv(envFile, env)
		case "tfvars", "hcl", "tfvar":
			err = service.WriteHcl(envFile, env)
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

	exportCmd.Flags().String("section", "", "The section name")

	exportCmd.Flags().String("format", "env", "The file format to save as (env, json, tfvars, hcl)")

}
