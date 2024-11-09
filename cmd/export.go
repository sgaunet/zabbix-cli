package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/spf13/cobra"
)

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

		z, err := zabbix.New(conf.ZabbixUser, conf.ZabbixPassword, conf.ZabbixEndpoint)
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
			res, err := z.Export(ctx, zabbix.ExportRequestOptionYAMLFormat(), zabbix.ExportRequestOptionTemplatesID([]string{tmplID}))
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err.Error())
				os.Exit(1)
			}
			fmt.Println(res)
		}
	},
}
