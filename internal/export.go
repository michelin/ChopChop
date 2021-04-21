package internal

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/jedib0t/go-pretty/table"
)

// ExporterFunc is a func type exporting the results
// to a writer.
type ExporterFunc func([]*Result, io.WriteCloser) error

var exportersMap = map[string]struct {
	ExporterFunc   ExporterFunc
	WriterProvider func(filename string) (io.WriteCloser, error)
}{
	"csv": {
		ExporterFunc: ExportCSV,
		WriterProvider: func(filename string) (io.WriteCloser, error) {
			return os.OpenFile(filename+".csv", os.O_CREATE|os.O_WRONLY, 0644)
		},
	},
	"json": {
		ExporterFunc: ExportJSON,
		WriterProvider: func(filename string) (io.WriteCloser, error) {
			return os.OpenFile(filename+".json", os.O_CREATE|os.O_WRONLY, 0644)
		},
	},
	"stdout": {
		ExporterFunc: ExportTable,
		WriterProvider: func(filename string) (io.WriteCloser, error) {
			return os.Stdout, nil
		},
	},
}

func ExportersList() string {
	s := ""
	l := len(exportersMap)
	i := 0
	for exp := range exportersMap {
		s += exp
		if i < l-1 {
			s += ", "
		}
		i++
	}
	return s
}

// ErrEmptyResults is an error meaning the results are empty.
var ErrEmptyResults = errors.New("no result found")

// ErrMaxSeverityReached is an error meaning a result severity
// has been reached.
type ErrMaxSeverityReached struct {
	Max, Sev Severity
}

func (e ErrMaxSeverityReached) Error() string {
	maxStr, _ := e.Max.String()
	sevStr, _ := e.Sev.String()
	return "max severity (" + maxStr + ") reached (" + sevStr + ")"
}

func CheckSeverities(results []*Result, max Severity) error {
	_, err := max.String()
	if err != nil {
		return err
	}

	// Check results severities are valid
	for _, res := range results {
		sevRes, err := StringToSeverity(res.Severity)
		if err != nil {
			return err
		}
		if sevRes < max {
			return &ErrMaxSeverityReached{Max: max, Sev: sevRes}
		}
	}
	return nil
}

// ErrUnsupportedExporter is an error meaning an exporter in
// a Config is not supported.
type ErrUnsupportedExporter struct {
	Exporter string
}

func (e ErrUnsupportedExporter) Error() string {
	return "unsupported exporter: " + e.Exporter
}

// ExportResults exports the results given a config, to a filename
// if the exporter needs it.
func ExportResults(results []*Result, config *Config, filename string) error {
	// Check parameters
	if config == nil {
		return &ErrNilParameter{"config"}
	}
	if len(results) == 0 {
		return ErrEmptyResults
	}

	// Check severities
	err := CheckSeverities(results, config.MaxSeverity)
	if err != nil {
		return err
	}

	// Export results
	exported := make(map[string]struct{})
	for _, format := range config.ExportFormats {
		if _, ok := exported[format]; !ok {
			exported[format] = struct{}{}

			d := exportersMap[format]
			f, err := d.WriterProvider(filename)
			if err != nil {
				return err
			}
			defer f.Close()

			err = d.ExporterFunc(results, f)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ExportJSON(results []*Result, w io.WriteCloser) error {
	// Marshal results in JSON
	jsonbytes, err := json.Marshal(results)
	if err != nil {
		return err
	}

	// Write JSON content in file
	if _, err := w.Write(jsonbytes); err != nil {
		return err
	}

	return nil
}

func ExportCSV(results []*Result, w io.WriteCloser) error {
	// Write headers
	_, err := w.Write([]byte("url,endpoint,severity,checkName,remediation\n"))
	if err != nil {
		return err
	}

	// Write content
	for _, result := range results {
		entry := result.URL + "," + result.Endpoint + "," + result.Severity + "," + result.Name + "," + result.Remediation + "\n"
		_, err := w.Write([]byte(entry))
		if err != nil {
			return err
		}
	}

	return nil
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

func ExportTable(results []*Result, w io.WriteCloser) error {
	t := table.NewWriter()
	t.SetOutputMirror(w)
	t.AppendHeader(table.Row{"URL", "Endpoint", "Severity", "Plugin", "Remediation"})
	for _, result := range results {
		// Convert and check severity
		sev, err := StringToSeverity(result.Severity)
		if err != nil {
			return err
		}

		// Build log severity
		var severity string
		switch sev {
		case High:
			severity += colorRed + "High"
		case Medium:
			severity += colorYellow + "Medium"
		case Low:
			severity += colorGreen + "Low"
		case Informational:
			severity += colorCyan + "Informational"
		}
		severity += colorReset

		// Append the content row
		t.AppendRow([]interface{}{
			result.URL,
			result.Endpoint,
			severity,
			result.Name,
			result.Remediation,
		})
	}
	t.SortBy([]table.SortBy{
		{Name: "Severity", Mode: table.Asc},
	})
	t.Render()

	return nil
}
