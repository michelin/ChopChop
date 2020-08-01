package app

import (
	"gochopchop/data"
	"io/ioutil"
	"log"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// List checks of config file
func List(cmd *cobra.Command, args []string) {

	configFile, _ := cmd.Flags().GetString("config-file")
	severity, _ := cmd.Flags().GetString("severity")

	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	fileData, err := ioutil.ReadAll(file)
	y := data.Config{}

	err = yaml.Unmarshal([]byte(fileData), &y)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	cpt := 0
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"URL", "Plugin Name", "Severity", "Description"})
	for index, plugin := range y.Plugins {
		_ = index
		for index, check := range plugin.Checks {
			_ = index
			// If the user wants a specific severity, collect only specified severity checks
			if severity != "" {
				if severity == string(*check.Severity) {
					t.AppendRow([]interface{}{plugin.URI, check.PluginName, *check.Severity, *check.Description})
					cpt++
				}
			} else {
				t.AppendRow([]interface{}{plugin.URI, check.PluginName, *check.Severity, *check.Description})
				cpt++
			}
		}
	}
	t.AppendFooter(table.Row{"", "", "Total Checks", cpt})
	t.Render()
}
