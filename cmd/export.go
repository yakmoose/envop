package cmd

import (
	"envop/service"
	"github.com/1Password/connect-sdk-go/connect"
	"github.com/spf13/cobra"
	"strings"
)

// importCmd represents the import command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export the specified 1password item into an environment file",
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

		item, err := service.Get1PasswordItem(client, vaultName, itemName)
		if err != nil {
			return err
		}

		fields := make(map[string]string, 0)
		for _, v := range item.Fields {
			fields[strings.ToUpper(v.Label)] = v.Value
		}

		format, err := cmd.Flags().GetString("format")
		if err != nil {
			return err
		}

		switch format {
		case "json":
			service.WriteJSON(environment, path, fields)
		case "env":
			service.WriteEnv(environment, path, fields)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().String("path", ".env", "The env file base")
	exportCmd.Flags().String("env", "dev", "The env environment")

	exportCmd.Flags().StringP("vault", "V", "", "The 1password vault")
	exportCmd.MarkFlagRequired("vault")

	exportCmd.Flags().StringP("item", "i", "", "The name of the item to save")
	exportCmd.MarkFlagRequired("item")

	exportCmd.Flags().String("format", "env", "The name of the item to save")

}
