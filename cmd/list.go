package cmd

import (
	"fmt"
	"os"

	"github.com/michelin/gochopchop/core"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

type listOptions struct {
	Severity string
}

func init() {
	pluginCmd := &cobra.Command{
		Use:   "plugins",
		Short: "list checks of configuration file",
		RunE:  runList,
	}
	addSignaturesFlag(pluginCmd)
	pluginCmd.Flags().StringP("severity", "s", "", "severity option for list tag") // --severity ou -s

	rootCmd.AddCommand(pluginCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	signatures, err := parseSignatures(cmd)
	if err != nil {
		return err
	}
	options, err := parseOptions(cmd)
	if err != nil {
		return err
	}
	cpt := 0
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"URL", "Plugin Name", "Severity", "Description"})
	for _, plugin := range signatures.Plugins {
		for _, check := range plugin.Checks {
			if options.Severity == "" || options.Severity == string(check.Severity) {
				t.AppendRow([]interface{}{plugin.Endpoint, check.Name, check.Severity, check.Description})
				cpt++
			}
		}
	}
	t.AppendFooter(table.Row{"", "", "Total Checks", cpt})
	t.Render()
	return nil
}

func parseOptions(cmd *cobra.Command) (*listOptions, error) {
	options := new(listOptions)
	severity, err := cmd.Flags().GetString("severity")
	if err != nil {
		return nil, fmt.Errorf("invalid value for severity: %v", err)
	}
	if severity != "" {
		if !core.ValidSeverity(severity) {
			return nil, fmt.Errorf("Invalid severity level : %s. Please use : %s", severity, core.SeveritiesAsString())
		}
		options.Severity = severity
	}
	return options, nil
}
