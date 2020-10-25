package formatting

import (
	"fmt"
	"gochopchop/core"
	"os"

	log "github.com/sirupsen/logrus"
)

// ExportCSV is a simple wrapper for CSV formatting
func ExportCSV(exportFilename string, out []core.Output) error {

	exportFilename = fmt.Sprintf("%s.csv", exportFilename)

	f, err := os.OpenFile(exportFilename, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("url,endpoint,severity,checkName,remediation\n")
	if err != nil {
		return err
	}
	for _, output := range out {
		line := fmt.Sprintf("%s,%s,%s,%s,%s\n", output.URL, output.Endpoint, output.Severity, output.Name, output.Remediation)
		_, err = f.Write([]byte(line))
		if err != nil {
			return err
		}
	}

	log.Info("Export as csv: ", exportFilename)
	return nil
}
