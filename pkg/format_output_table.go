package pkg

import (
	"fmt"
	"gochopchop/data"
	"os"

	"github.com/jedib0t/go-pretty/table"
)

// FormatOutputTable will render the data as a nice table
func FormatOutputTable(out []data.Output) {

	colorReset := "\033[0m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	colorYellow := "\033[33m"
	// colorBlue := "\033[34m"
	// colorPurple := "\033[35m"
	colorCyan := "\033[36m"
	// colorWhite := "\033[37m"

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Domain", "URI", "Severity", "Plugin", "Remediation"})
	for i := 0; i < len(out); i++ {
		severity := ""
		if out[i].Severity == "High" {
			severity = fmt.Sprint(string(colorRed), "High", string(colorReset))
		} else if out[i].Severity == "Medium" {
			severity = fmt.Sprint(string(colorYellow), "Medium", string(colorReset))
		} else if out[i].Severity == "Low" {
			severity = fmt.Sprint(string(colorGreen), "Low", string(colorReset))
		} else {
			severity = fmt.Sprint(string(colorCyan), "Informational", string(colorReset))
		}

		t.AppendRow([]interface{}{
			out[i].Domain,
			out[i].TestedURL,
			severity,
			out[i].PluginName,
			out[i].Remediation,
		})
	}
	t.SortBy([]table.SortBy{
		{Name: "Severity", Mode: table.Asc},
	})
	t.Render()
}
