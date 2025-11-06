package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import template",
	Long:  `import template`,
	Run: func(_ *cobra.Command, _ []string) {
		var err error
		ctx := context.Background()

		err = initConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}

		if templateFile == "" {
			fmt.Fprintf(os.Stderr, "template file is mandatory\n")
			os.Exit(1)
		}

		z := zabbix.New(conf.ZabbixUser, conf.ZabbixPassword, conf.ZabbixEndpoint)
		err = z.Login(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}
		defer z.Logout(ctx) //nolint: errcheck

		// read file
		template, err := os.ReadFile(templateFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}

		templateData := string(template)

		// Detect format from file extension
		detectedFormat := zabbix.DetectFormatFromExtension(templateFile)
		if detectedFormat == zabbix.FormatUnknown {
			// Try to detect from content
			detectedFormat = zabbix.DetectFormatFromContent(templateData)
		}

		// Validate format
		if detectedFormat != zabbix.FormatUnknown {
			if err := zabbix.ValidateFormat(templateData, detectedFormat); err != nil {
				fmt.Fprintf(os.Stderr, "Format validation failed: %v\n", err.Error())
				os.Exit(1)
			}
		}

		// Validate it's Zabbix export data
		if err := zabbix.ValidateZabbixExportData(templateData); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid Zabbix export data: %v\n", err.Error())
			os.Exit(1)
		}

		// Import the template
		isImported, err := z.Import(ctx, templateData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}
		if isImported {
			fmt.Println("Template imported")
		} else {
			fmt.Println("Template not imported")
		}
	},
}
