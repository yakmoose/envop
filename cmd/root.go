/*
Copyright © 2025 John Lennard <john@yakmoo.se>
*/
package cmd

import (
	"github.com/spf13/pflag"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "envop",
	Short: "Imports environment files into 1password",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.envop.json)")
	rootCmd.PersistentFlags().StringP("service-account", "", "", "1password service account")
	rootCmd.MarkPersistentFlagRequired("service-account")

	viper.BindPFlag("service-account", rootCmd.PersistentFlags().Lookup("service-account"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".envop.json")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer("SERVICE-ACCOUNT", "OP_SERVICE_ACCOUNT_TOKEN"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// panic(err)
	}

	rootCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		// Determine the naming convention of the flags when represented in the config file
		configName := f.Name

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(configName) {
			rootCmd.PersistentFlags().Set(f.Name, viper.GetString(configName))
		}
	})
}
