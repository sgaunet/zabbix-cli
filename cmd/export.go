package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/spf13/cobra"
)

// GetFormatOption returns the appropriate export format option based on the format string
func GetFormatOption(format string) (zabbix.ConfigurationExportRequestOption, error) {
	switch strings.ToLower(format) {
	case "yaml":
		return zabbix.ExportRequestOptionYAMLFormat(), nil
	case "json":
		return zabbix.ExportRequestOptionJSONFormat(), nil
	case "xml":
		return zabbix.ExportRequestOptionXMLFormat(), nil
	default:
		return nil, fmt.Errorf("invalid format: %s (valid formats: yaml, json, xml)", format)
	}
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export template",
	Long:  `export template`,
	Run: func(_ *cobra.Command, _ []string) {
		var err error
		ctx := context.Background()

		err = initConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}

		if templateName == "" {
			fmt.Fprintf(os.Stderr, "template name is required\n")
			os.Exit(1)
		}

		// Get format option
		formatOpt, err := GetFormatOption(exportFormat)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}

		z := zabbix.New(conf.ZabbixUser, conf.ZabbixPassword, conf.ZabbixEndpoint)
		err = z.Login(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}
		defer z.Logout(ctx) //nolint: errcheck

		res, err := z.GetTemplates(zabbix.GetTemplateOptionFilterByName([]string{templateName}))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}
		templatesID := res.GetTemplateID()
		for _, tmplID := range templatesID {
			// fmt.Println(template.TemplateID, template.Name)
			res, err := z.Export(ctx, formatOpt, zabbix.ExportRequestOptionTemplatesID([]string{tmplID}))
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err.Error())
				os.Exit(1)
			}
			fmt.Println(res)
		}
	},
}
