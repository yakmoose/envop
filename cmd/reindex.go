/*
Copyright © 2025 John Lennard <john@yakmoo.se>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yakmoose/envop/service"
)

// reindexCmd reindex the thing
var reindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Reindex the specified item",
	RunE: func(cmd *cobra.Command, args []string) error {

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

		client, err := service.NewClientFromToken(token)
		if err != nil {
			return err
		}

		vault, err := service.FindVaultWithName(client, vaultName)
		if err != nil {
			return err
		}

		item, err := service.FindItemWithName(client, vault, itemName)
		if err != nil {
			return err
		}

		if item == nil {
			return fmt.Errorf("Item %s not found in vault %s", itemName, vaultName)
		}

		_, err = service.ReindexItem(client, item)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(reindexCmd)

	reindexCmd.Flags().String("vault", "", "The 1password vault")
	reindexCmd.MarkFlagRequired("vault")

	reindexCmd.Flags().String("item", "", "The name of the item to save")
	reindexCmd.MarkFlagRequired("item")

}
