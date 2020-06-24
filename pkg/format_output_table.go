package pkg

import (
	"gochopchop/data"
	"os"

	"github.com/jedib0t/go-pretty/table"
)

// FormatOutputTable will render the data as a nice table
func FormatOutputTable(out []data.Output) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Domain", "URL", "Severity", "Plugin", "Remediation"})
	for i := 0; i < len(out); i++ {
		t.AppendRow([]interface{}{
			out[i].Domain,
			out[i].TestedURL,
			out[i].Severity,
			out[i].PluginName,
			out[i].Remediation,
		})
	}
	t.Render()
}
