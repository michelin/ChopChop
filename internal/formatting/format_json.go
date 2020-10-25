package formatting

import (
	"encoding/json"
	"fmt"
	"gochopchop/core"
	"os"

	log "github.com/sirupsen/logrus"
)

// ExportJSON will save the output to a JSON file
func ExportJSON(exportFilename string, output []core.Output) error {

	jsonstr, err := json.Marshal(output)
	if err != nil {
		return err
	}

	exportFilename = fmt.Sprintf("./%s.json", exportFilename)

	f, err := os.OpenFile(exportFilename, os.O_CREATE|os.O_WRONLY, 0755)
	defer f.Close()
	if err != nil {
		return err
	}

	_, err = f.Write(jsonstr)
	if err != nil {
		return err
	}

	log.Info("Export as json: ", exportFilename)

	return nil
}
