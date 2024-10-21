package cmd

import (
	"fmt"
	"os"

	"github.com/sgaunet/zabbix-cli/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string = "default"
var templateName string
var conf *config.Config

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
	// export subcommand
	exportCmd.Flags().StringVar(&templateName, "t", "", "template name to export")
	rootCmd.AddCommand(exportCmd)
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
		fmt.Fprintf(os.Stderr, "unable to decode into config struct, %v", err)
		os.Exit(1)
	}
	if !conf.IsValid() {
		return fmt.Errorf("config is not valid")
	}
	return nil
}
