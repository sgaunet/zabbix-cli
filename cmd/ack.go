package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/spf13/cobra"
)

var ackCmd = &cobra.Command{
	Use:   "ack",
	Short: "acknowledge events",
	Long:  `acknowledge events`,
	Run: func(_ *cobra.Command, _ []string) {
		var err error
		ctx := context.Background()

		err = initConfig()
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

		problems, err := z.GetProblems(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}

		for _, pb := range problems {
			fmt.Printf("Acknowledge problem %s\n", pb.Name)
			_, err = z.AcknowledgeEvents(ctx, []string{pb.EventID}, zabbix.WithActions(zabbix.AddMessage, zabbix.CloseProblem, zabbix.Acknowledge), zabbix.WithMessage("acknowledged from CLI"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err.Error())
				os.Exit(1)
			}
		}
	},
}
