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

var Ack bool
var Supp bool
var problemSeverityFlag string

// ProblemGetCmd represents the get problem subcommand
var ProblemGetCmd = &cobra.Command{
	Use:   "get",
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

		var options []zabbix.GetProblemOption
		options = append(options, zabbix.GetProblemOptionAcknowledged(Ack))
		options = append(options, zabbix.GetProblemOptionSuppressed(Supp))

		if problemSeverityFlag != "" {
			severityInt := zabbix.GetSeverityString(problemSeverityFlag)
			// ProblemParams.Severities expects []string of integer severities
			options = append(options, zabbix.GetProblemOptionSeverities([]string{fmt.Sprintf("%d", severityInt)}))
		}
		// Add SelectHosts to get host information
		options = append(options, zabbix.GetProblemOptionSelectHosts("extend"))

		res, err := z.GetProblems(ctx, options...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}

		PrettyPrintProblems(res) //nolint:errcheck
	},
}

// getSeverityStyle returns the appropriate pterm style for a given severity level
func getSeverityStyle(severity string) *pterm.Style {
	switch severity {
	case "Information":
		return pterm.NewStyle(pterm.FgBlue)
	case "Warning":
		return pterm.NewStyle(pterm.FgYellow)
	case "Average":
		return pterm.NewStyle(pterm.FgLightYellow)
	case "High":
		return pterm.NewStyle(pterm.FgRed)
	case "Disaster":
		return pterm.NewStyle(pterm.FgRed, pterm.Bold)
	default: // Not classified
		return pterm.NewStyle(pterm.FgWhite)
	}
}

func PrettyPrintProblems(problems []zabbix.Problem) error {
	tData := pterm.TableData{
		// Header row for the table
		{"Time", "Host", "Problem", "Severity", "Ack", "Suppressed", "Duration"},
	}

	for _, pb := range problems {
		severity := pb.GetSeverity()
		severityStyle := getSeverityStyle(severity)

		// Format the severity with appropriate color
		coloredSeverity := severityStyle.Sprint(severity)

		var hostName string
		if len(pb.Hosts) > 0 {
			hostName = pb.Hosts[0].Name // Assuming the first host is the relevant one
		} else {
			hostName = "N/A"
		}

		tData = append(tData, []string{
			pb.GetClock().Format("2006-01-02 15:04:05"), // Time
			hostName,                                   // Host
			pb.Name,                                    // Problem
			coloredSeverity,                            // Severity
			pb.GetAcknowledgeStr(),                     // Ack
			pb.GetSuppressedStr(),                      // Suppressed
			pb.GetDurationStr(),                        // Duration
		})
	}

	// Create a table with a header and the defined data, then render it
	err := pterm.DefaultTable.
		WithHasHeader().
		WithBoxed(true).
		WithData(tData).
		Render()

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
