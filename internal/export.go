package internal

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/table"
	log "github.com/sirupsen/logrus"
)

// ExporterFunc is a func type exporting the results
// to a file given its filename.
type ExporterFunc func([]*Result, string) error

var exportersMap = map[string]ExporterFunc{
	"csv":    exportCSV,
	"json":   exportJSON,
	"stdout": exportStdout,
}

var ErrEmptyResults = errors.New("no result found")
var ErrNilConfig = errors.New("given config is nil")
var ErrMaxSeverityReached = errors.New("max severity reached")

func ExportResults(results []*Result, config *Config, filename string) error {
	// Check parameters
	if config == nil {
		return ErrNilConfig
	}
	if len(results) == 0 {
		return ErrEmptyResults
	}

	// Check results severities are valid
	maxSevStr := config.MaxSeverity
	for _, res := range results {
		sevRes, err := StringToSeverity(res.Severity)
		if err != nil {
			return err
		}
		if sevRes > maxSevStr {
			return ErrMaxSeverityReached
		}
	}

	// Export results
	exported := make(map[string]struct{})
	for _, format := range config.ExportFormats {
		if _, ok := exported[format]; !ok {
			exported[format] = struct{}{}

			f, ok := exportersMap[format]
			if !ok {
				return errors.New("unsupported exporter")
			}
			if err := f(results, filename); err != nil {
				return err
			}
		}
	}

	return nil
}

// exportJSON exports the results to a JSON file
func exportJSON(results []*Result, filename string) error {
	// Open CSV file to write in
	f, err := os.OpenFile(filename+".json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Marshal results in JSON
	jsonbytes, err := json.Marshal(results)
	if err != nil {
		return err
	}

	// Write JSON content in file
	if _, err := f.Write(jsonbytes); err != nil {
		return err
	}

	log.Info("Results were exported as json in: ", filename)
	return nil
}

// exportCSV exports the results to a CSV file
func exportCSV(results []*Result, filename string) error {
	// Open CSV file to write in
	f, err := os.OpenFile(filename+".csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write headers
	_, err = f.WriteString("url,endpoint,severity,checkName,remediation\n")
	if err != nil {
		return err
	}

	// Write content
	for _, result := range results {
		entry := result.URL + "," + result.Endpoint + "," + result.Severity + "," + result.Name + "," + result.Remediation + "\n"
		_, err := f.WriteString(entry)
		if err != nil {
			return err
		}
	}

	log.Info("Results were exported as csv in: ", filename)
	return nil
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

type ErrUnsupportedSeverity struct {
	Severity Severity
}

func (e ErrUnsupportedSeverity) Error() string {
	return "unsupported severity " + strconv.Itoa(int(e.Severity))
}

// exportStdout export the results in os.Stdout
func exportStdout(results []*Result, filename string) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
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
		default:
			return &ErrUnsupportedSeverity{sev}
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
