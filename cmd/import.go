/*
Copyright © 2025 John Lennard <john@yakmoo.se>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yakmoose/envop/service"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import the specified file into 1password",
	RunE: func(cmd *cobra.Command, args []string) error {
		envFile, err := cmd.Flags().GetString("env-file")
		if err != nil {
			return err
		}

		envName, err := cmd.Flags().GetString("env-name")
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

		environment, err := service.ReadEnv(envName, envFile)
		if err != nil {
			return err
		}

		if len(environment) > 0 {

			vault, err := service.FindVaultWithName(client, vaultName)
			if err != nil {
				return err
			}

			item, err := service.FindItemWithName(client, vault, itemName)
			if err != nil {
				return err
			}

			replace, err := cmd.Flags().GetBool("replace")
			if err != nil {
				return err
			}

			if item != nil && replace {
				err := client.Items().Delete(context.Background(), vault.ID, item.ID)
				if err != nil {
					return err
				}
				item = nil
			}

			if item == nil {
				item, err = service.CreateItemInVaultWithSection(
					client,
					vault,
					itemName,
					sectionName,
					&environment,
				)
			} else {
				item, err = service.UpdateItem(
					client,
					item,
					sectionName,
					&environment,
				)
			}

			if err != nil {
				return err
			}
			fmt.Printf("item created: %s (%s)\n", item.Title, item.ID)

		} else {
			return fmt.Errorf("no items found for environment: %s", envName)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().String("env-file", "", "The env file base")
	importCmd.Flags().String("env-name", "", "The environment, will try <path>, <path>.local, <path>.<env> and <path>.<env>.local")

	importCmd.Flags().String("vault", "", "The 1password vault")
	importCmd.MarkFlagRequired("vault")

	importCmd.Flags().String("section", "", "The 1password section to add fields to")
	importCmd.MarkFlagRequired("section")

	importCmd.Flags().String("item", "", "The name of the item to save")
	importCmd.MarkFlagRequired("item")

	importCmd.Flags().Bool("replace", false, "Replace existing item instead of appending to it")
}
