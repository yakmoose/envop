/*
Copyright © 2023 John Lennard <john@yakmoo.se>
*/
package cmd

import (
	"fmt"

	"envop/service"

	"github.com/1Password/connect-sdk-go/connect"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import the specified file into 1password",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return err
		}

		environment, err := cmd.Flags().GetString("env")
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

		client, err := connect.NewClientFromEnvironment()
		if err != nil {
			return err
		}

		item, err := service.Create1PasswordItem(
			client,
			vaultName,
			itemName,
			service.ReadEnv(environment, path),
		)
		if err != nil {
			return err
		}
		fmt.Printf("item created: %s (%s)\n", item.Title, item.ID)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().String("path", ".env", "The env file base")
	importCmd.Flags().String("env", "dev", "The env environment")

	importCmd.Flags().StringP("vault", "V", "", "The 1password vault")
	importCmd.MarkFlagRequired("vault")

	importCmd.Flags().String("item", "i", "The name of the item to save")
	importCmd.MarkFlagRequired("item")
}
