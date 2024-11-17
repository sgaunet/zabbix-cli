package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/sgaunet/zabbix-cli/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var conf *config.Config // configuration
var cfgFile string      // permanent flag to specify configuration file

// templateName is a flag to specify the template name (export)
var templateName string

// templateFile is a flag to specify the template file (import)
var templateFile string

var ErrInvalidConfig = errors.New("invalid configuration")

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zabbix-cli",
	Short: "A CLI tool to interact with Zabbix API",
	Long:  `A CLI tool to interact with Zabbix API`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "default", "configuration file (default is $HOME/.config/zabbix-cli/default.yaml)")

	// export subcommand
	exportCmd.Flags().StringVarP(&templateName, "template", "t", "", "template name to export")
	// import subcommand
	importCmd.Flags().StringVarP(&templateFile, "file", "f", "", "template file to import")
	importCmd.MarkFlagRequired("f") //nolint: errcheck

	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)

	GetCmd.AddCommand(GetProblemCmd)
	rootCmd.AddCommand(GetCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot get user home directory: %w", err)
	}

	// Search config in home directory/.config/zabbix-cli
	configDir := fmt.Sprintf("%s/%s", home, ".config/zabbix-cli")
	viper.AddConfigPath(configDir)
	viper.SetConfigType("yaml")
	viper.SetConfigName(cfgFile)
	_ = viper.ReadInConfig() // no problem if env variables are set and there is a check after

	_ = viper.BindEnv("ZABBIX_ENDPOINT")
	_ = viper.BindEnv("ZABBIX_USER")
	_ = viper.BindEnv("ZABBIX_PASSWORD")
	viper.AutomaticEnv()

	conf = &config.Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		return fmt.Errorf("unable to decode into config struct, %w", err)
	}
	if !conf.IsValid() {
		return ErrInvalidConfig
	}
	return nil
}
