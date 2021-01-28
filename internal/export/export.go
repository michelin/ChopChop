package export

import (
	"encoding/json"
	"fmt"
	"gochopchop/core"
	"os"

	log "github.com/sirupsen/logrus"
)

type IFile interface {
	WriteString(input string) (n int, err error)
}

// ExportCSV exports the output in a CSV file
func ExportCSV(filename string, out []core.Output) error {
	exportFilename := fmt.Sprintf("%s.csv", filename)

	f, err := os.OpenFile(exportFilename, os.O_CREATE|os.O_WRONLY, 0755)
	defer f.Close()
	if err != nil {
		return err
	}

	err = exportCSV(f, out)
	if err != nil {
		return err
	}
	log.Info("Results were exported as csv in: ", exportFilename)
	return nil
}

func exportCSV(file IFile, out []core.Output) error {
	_, err := file.WriteString("url,endpoint,severity,checkName,remediation\n")
	if err != nil {
		return err
	}
	for _, output := range out {
		line := fmt.Sprintf("%s,%s,%s,%s,%s\n", output.URL, output.Endpoint, output.Severity, output.Name, output.Remediation)
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}

	return nil
}

// ExportJSON will save the output to a JSON file
func ExportJSON(filename string, output []core.Output) error {
	exportFilename := fmt.Sprintf("%s.json", filename)

	f, err := os.OpenFile(exportFilename, os.O_CREATE|os.O_WRONLY, 0755)
	defer f.Close()
	if err != nil {
		return err
	}

	err = exportJSON(f, output)
	if err != nil {
		return err
	}
	log.Info("Results were exported as json in: ", exportFilename)
	return nil
}

func exportJSON(file IFile, output []core.Output) error {
	jsonbytes, err := json.Marshal(output)
	if err != nil {
		return err
	}
	if _, err := file.WriteString(string(jsonbytes)); err != nil {
		return err
	}
	return nil
}
