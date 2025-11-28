/*
Copyright © 2025 John Lennard <john@yakmoo.se>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yakmoose/envop/service"
)

// copyCmd copy a section from one item to another
var copyCmd = &cobra.Command{
	Use:   "cp",
	Short: "Copy the specified section to a new item",
	RunE: func(cmd *cobra.Command, args []string) error {

		sourceVaultName, err := cmd.Flags().GetString("source-vault")
		if err != nil {
			return err
		}

		sourceItemName, err := cmd.Flags().GetString("source-item")
		if err != nil {
			return err
		}

		destinationVaultName, _ := cmd.Flags().GetString("destination-vault")
		if destinationVaultName == "" {
			destinationVaultName = sourceVaultName
		}

		destinationItemName, _ := cmd.Flags().GetString("destination-item")
		if destinationItemName == "" {
			destinationItemName = sourceItemName
		}

		sourceSectionName, _ := cmd.Flags().GetString("source-section")
		if err != nil {
			return err
		}

		destinationSectionName, _ := cmd.Flags().GetString("destination-section")
		if destinationSectionName == "" {
			destinationSectionName = sourceSectionName
		}

		token, err := cmd.Flags().GetString("service-account")
		if err != nil {
			return err
		}

		client, err := service.NewClientFromToken(token)
		if err != nil {
			return err
		}

		sourceVault, err := service.FindVaultWithName(client, sourceVaultName)
		if err != nil {
			return err
		}

		sourceItem, err := service.FindItemWithName(client, sourceVault, sourceItemName)
		if err != nil {
			return err
		}

		if sourceItem == nil {
			return fmt.Errorf("Item %s not found in vault %s", sourceItemName, sourceItemName)
		}

		destinationVault, err := service.FindVaultWithName(client, destinationVaultName)
		if err != nil {
			return err
		}

		destinationItem, err := service.FindItemWithName(client, destinationVault, destinationItemName)
		if err != nil {
			return err
		}

		return service.CopySection(
			client,
			sourceItem,
			sourceSectionName,
			destinationItem,
			destinationSectionName,
		)
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	copyCmd.Flags().String("source-vault", "", "The 1password vault to copy from")
	copyCmd.MarkFlagRequired("source-vault")

	copyCmd.Flags().String("source-item", "", "The name of the item to copy from")
	copyCmd.MarkFlagRequired("source-item")

	copyCmd.Flags().String("source-section", "", "The 1password section to copy from ")
	copyCmd.MarkFlagRequired("source-section")

	copyCmd.Flags().String("destination-vault", "", "The 1password vault to copy to")
	copyCmd.Flags().String("destination-item", "", "The name of the item to copy to")
	copyCmd.Flags().String("destination-section", "", "The 1password section to copy to")
}
