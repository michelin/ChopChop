package cmd

import (
	"fmt"
	"gochopchop/core"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var signatureFlagName = "signatures"
var signatureFlagShorthand = "c"
var signatureDefaultFilename = "chopchop.yml"

func addSignaturesFlag(cmd *cobra.Command) error {
	cmd.Flags().StringP(signatureFlagName, signatureFlagShorthand, signatureDefaultFilename, "path to signature file") // --signature ou -c
	return nil
}

func parseSignatures(cmd *cobra.Command) (*core.Signatures, error) {

	signatureFile, err := cmd.Flags().GetString(signatureFlagName)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for signatureFile: %v", err)
	}
	if _, err := os.Stat(signatureFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("Path of signatures file is not valid")
	}

	file, err := os.Open(signatureFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	signatureData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	signatures := core.NewSignatures()

	err = yaml.Unmarshal([]byte(signatureData), signatures)
	if err != nil {
		return nil, err
	}

	severityFilter, _ := cmd.Flags().GetString("severity-filter")
	if severityFilter != "" {
		signatures.FilterBySeverity(severityFilter)
	}

	pluginFilters, _ := cmd.Flags().GetStringSlice("plugin-filters")
	if len(pluginFilters) > 0 {
		signatures.FilterByNames(pluginFilters)
	}

	for _, plugin := range signatures.Plugins {
		if plugin.Endpoint == "" {
			if len(plugin.Endpoints) > 0 {
				return nil, fmt.Errorf("URI and URIs can't be set at the same time in plugin checks. Stopping execution.")
			}
		}
		for _, check := range plugin.Checks {
			if check.Description == "" {
				return nil, fmt.Errorf("Missing or empty description field in %s plugin checks. Stopping execution.", check.Name)
			}
			if check.Remediation == "" {
				return nil, fmt.Errorf("Missing or empty remediation field in %s plugin checks. Stopping execution.", check.Name)
			}
			if check.Severity == "" {
				return nil, fmt.Errorf("Missing severity field in %s plugin checks. Stopping execution.", check.Name)
			}
			if !core.ValidSeverity(check.Severity) {
				return nil, fmt.Errorf("Invalid severity : %s. Please use : %s", check.Severity, core.SeveritiesAsString())
			}
		}
	}

	return signatures, nil
}
