/*
Copyright © 2025 John Lennard <john@yakmoo.se>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yakmoose/envop/service"
)

// moveCmd move a section from one item to another
var moveCmd = &cobra.Command{
	Use:   "mv",
	Short: "Move the specified section to a new item",
	RunE: func(cmd *cobra.Command, args []string) error {

		sourceVaultName, err := cmd.Flags().GetString("source-vault")
		if err != nil {
			return err
		}

		sourceItemName, err := cmd.Flags().GetString("source-item")
		if err != nil {
			return err
		}

		sourceSectionName, _ := cmd.Flags().GetString("source-section")
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
			return fmt.Errorf("Item %s not found in vault %s", sourceItemName, sourceVaultName)
		}

		destinationVault, err := service.FindVaultWithName(client, destinationVaultName)
		if err != nil {
			return err
		}

		destinationItem, err := service.FindItemWithName(client, destinationVault, destinationItemName)
		if err != nil {
			return err
		}

		return service.MoveSection(
			client,
			sourceItem,
			sourceSectionName,
			destinationItem,
			destinationSectionName,
		)
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)

	moveCmd.Flags().String("source-vault", "", "The 1password vault to move from")
	moveCmd.MarkFlagRequired("source-vault")
	moveCmd.Flags().String("source-item", "", "The name of the item to move from")
	moveCmd.MarkFlagRequired("source-item")
	moveCmd.Flags().String("source-section", "", "The 1password section to move from ")
	moveCmd.MarkFlagRequired("source-section")

	moveCmd.Flags().String("destination-vault", "", "The 1password vault to move to")
	moveCmd.Flags().String("destination-item", "", "The name of the item to move to")
	moveCmd.Flags().String("destination-section", "", "The 1password section to move to")

}
