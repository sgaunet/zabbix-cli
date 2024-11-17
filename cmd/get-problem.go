package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/spf13/cobra"
)

// GetProblemCmd represents the get problem subcommand
var GetProblemCmd = &cobra.Command{
	Use:   "problem",
	Short: "get problems",
	Long:  `get problems`,
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
		defer z.Logout(ctx) //nolint:errcheck

		res, err := z.GetProblems(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}

		PrettyPrintProblems(res) //nolint:errcheck
	},
}

func PrettyPrintProblems(problems []zabbix.Problem) error {
	tData := pterm.TableData{
		{"Time", "Probem", "Severity", "Duration", "Ack", "Supp"},
	}
	for _, pb := range problems {
		tData = append(tData, []string{pb.GetClock().Format("2006-01-02 15:04:05"),
			pb.Name, pb.GetSeverity(), pb.GetDurationStr(), pb.GetAcknowledgeStr(), pb.GetSuppressedStr()})
	}
	// Create a table with a header and the defined data, then render it
	err := pterm.DefaultTable.WithHasHeader().WithData(tData).Render()
	if err != nil {
		return fmt.Errorf("error rendering table: %w", err)
	}
	return nil
}

func PrintProblemsCSV(problems []zabbix.Problem) error {
	csvwriter := csv.NewWriter(os.Stdout)
	err := csvwriter.Write([]string{"Time", "Probem", "Severity", "Duration", "Ack", "Supp"})
	if err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	for _, pb := range problems {
		err = csvwriter.Write([]string{pb.GetClock().Format("2006-01-02 15:04:05"),
			pb.Name, pb.GetSeverity(), pb.GetDurationStr(), pb.GetAcknowledgeStr(), pb.GetSuppressedStr()})
		if err != nil {
			return fmt.Errorf("error writing CSV row: %w", err)
		}
	}

	csvwriter.Flush()
	return nil
}
