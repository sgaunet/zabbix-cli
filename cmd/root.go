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
	cobra.OnInitialize(func() {
		if err := initConfig(); err != nil {
			// It's common to os.Exit(1) or panic here if config is essential and fails to load.
			// For now, let's print the error to stderr. Commands should check if conf is nil.
			fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
			// os.Exit(1) // Or handle more gracefully depending on desired behavior
		}
	})

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "default", "configuration file (default is $HOME/.config/zabbix-cli/default.yaml)")

	// export subcommand
	exportCmd.Flags().StringVarP(&templateName, "template", "t", "", "template name to export")
	// import subcommand
	importCmd.Flags().StringVarP(&templateFile, "file", "f", "", "template file to import")
	importCmd.MarkFlagRequired("f") //nolint: errcheck

	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)

	ProblemGetCmd.Flags().BoolVarP(&Ack, "ack", "a", false, "show acknowledged problems")
	ProblemGetCmd.Flags().BoolVarP(&Supp, "supp", "s", false, "show suppressed problems")
	ProblemCmd.AddCommand(ProblemGetCmd)
	rootCmd.AddCommand(ProblemCmd)

	rootCmd.AddCommand(ackCmd)
	rootCmd.AddCommand(MaintenanceCmd)
	MaintenanceCmd.AddCommand(MaintenanceCreateCmd)
	MaintenanceCmd.AddCommand(MaintenanceDeleteCmd)
	MaintenanceCmd.AddCommand(MaintenanceGetCmd)
	MaintenanceDeleteCmd.AddCommand(MaintenanceDeleteAllCmd)

	rootCmd.AddCommand(HostgroupCmd)
	HostgroupCmd.AddCommand(HostGroupGetCmd)
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
