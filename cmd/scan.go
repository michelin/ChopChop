package cmd

import (
	"bufio"
	"fmt"
	"gochopchop/core"
	"gochopchop/internal/export"
	"gochopchop/internal/formatting"
	"gochopchop/internal/httpget"
	"net/url"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "scan endpoints to check if services/files/folders are exposed",
		RunE:  runScan,
	}
	addSignaturesFlag(scanCmd)

	scanCmd.Flags().BoolP("insecure", "k", false, "Check SSL certificate")                                                                                    // --insecure ou -n
	scanCmd.Flags().StringP("url-file", "u", "", "path to a specified file containing urls to test")                                                          // --uri-file ou -f
	scanCmd.Flags().StringP("max-severity", "b", "", "block the CI pipeline if severity is over or equal specified flag")                                     // --max-severity ou -m
	scanCmd.Flags().StringSliceP("export", "e", []string{}, "export of the output (csv and json)")                                                            //--export ou --e
	scanCmd.Flags().StringP("export-filename", "", "", "filename for export files")                                                                           // --export-filename
	scanCmd.Flags().IntP("timeout", "t", 10, "Timeout for the HTTP requests (default: 10s)")                                                                  // --timeout ou -ts
	scanCmd.Flags().StringP("severity-filter", "", "", "Filter by severity (engine will check for same severity checks)")                                     // --severity-filter
	scanCmd.Flags().StringSliceP("plugin-filters", "", []string{}, "Filter by the name of the plugin (engine will only check for plugin with the same name)") // --plugin-filter
	rootCmd.AddCommand(scanCmd)
}

func runScan(cmd *cobra.Command, args []string) error {
	config, err := parseConfig(cmd, args)
	if err != nil {
		return err
	}

	signatures, err := parseSignatures(cmd)
	if err != nil {
		return err
	}

	begin := time.Now()

	fetcher := httpget.NewFetcher(config.HTTP.Insecure, config.HTTP.Timeout)
	noRedirectFetcher := httpget.NewNoRedirectFetcher(config.HTTP.Insecure, config.HTTP.Timeout)

	scanner := core.NewScanner(fetcher, noRedirectFetcher, signatures, config.Threads)

	result, err := scanner.Scan(cmd.Context(), config.Urls)
	if err != nil {
		return err
	}

	log.Info("Scan execution time:", time.Since(begin))

	if len(result) > 0 {

		formatting.PrintTable(result, os.Stdout)

		if contains(config.ExportFormats, "json") {
			export.ExportJSON(config.ExportFilename, result)
		}
		if contains(config.ExportFormats, "csv") {
			export.ExportCSV(config.ExportFilename, result)
		}

		if config.MaxSeverity != "" {
			for _, output := range result {
				if core.SeverityReached(config.MaxSeverity, output.Severity) {
					return fmt.Errorf("Max severity level reached, exiting with error code")
				}
			}
		}
	} else {
		log.Info("No vulnerabilities found. Exiting...")
	}
	return nil
}

func parseConfig(cmd *cobra.Command, args []string) (*core.Config, error) {

	urlFile, err := cmd.Flags().GetString("url-file")
	if err != nil {
		return nil, fmt.Errorf("invalid value for url-file: %v", err)
	}

	if urlFile != "" && len(args) >= 1 {
		// both urlFile and url are set, abort
		return nil, fmt.Errorf("Can't specify url with url list flag")
	}
	if urlFile == "" && len(args) == 0 {
		// no urlFile and no argument, abort
		return nil, fmt.Errorf("No url provided, please set the input-file flag or provide an url as an argument")
	}

	var urls []string
	if urlFile != "" {
		content, err := os.Open(urlFile)
		if err != nil {
			return nil, err
		}
		defer content.Close()
		scanner := bufio.NewScanner(content)
		for scanner.Scan() {
			url := scanner.Text()
			if !isURL(url) {
				log.Warn("url: ", url, " - is not valid - skipping scan")
				continue
			}
			urls = append(urls, url)
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	if len(args) > 1 {
		return nil, fmt.Errorf("Please provide only one URL")
	}

	if len(args) == 1 {
		url := args[0]
		if isURL(url) {
			urls = append(urls, url)
		} else {
			return nil, fmt.Errorf("Please provide a valid URL")
		}
	}

	insecure, err := cmd.Flags().GetBool("insecure")
	if err != nil {
		return nil, fmt.Errorf("invalid value for insecure: %v", err)
	}

	severityFilter, err := cmd.Flags().GetString("severity-filter")
	if err != nil {
		return nil, fmt.Errorf("invalid value for severity-filter: %v", err)
	}
	if severityFilter != "" {
		if !core.ValidSeverity(severityFilter) {
			return nil, fmt.Errorf("Invalid severity level : %s. Please use : %s", severityFilter, core.SeveritiesAsString())
		}
	}

	pluginFilters, err := cmd.Flags().GetStringSlice("plugin-filters")
	if err != nil {
		return nil, fmt.Errorf("invalid value for plugin-filters: %v", err)
	}

	exportFormats, err := cmd.Flags().GetStringSlice("export")
	if err != nil {
		return nil, fmt.Errorf("invalid value for export formats: %v", err)
	}
	if len(exportFormats) > 0 {
		for _, f := range exportFormats {
			if f != "csv" && f != "json" {
				return nil, fmt.Errorf("invalid value for export: %v , expected csv or json", f)
			}
		}
	}

	maxSeverity, err := cmd.Flags().GetString("max-severity")
	if err != nil {
		return nil, fmt.Errorf("invalid value for max sevirity : %v", err)
	}
	if maxSeverity != "" && !core.ValidSeverity(maxSeverity) {
		return nil, fmt.Errorf("Invalid max severity level : %s. Please use : %s", maxSeverity, core.SeveritiesAsString())
	}

	exportFilename, err := cmd.Flags().GetString("export-filename")
	if err != nil {
		return nil, fmt.Errorf("invalid value for exportFilename: %v", err)
	}
	if exportFilename == "" {
		now := time.Now().Format("2006-01-02_15-04-05")
		exportFilename = fmt.Sprintf("gochopchop_%s", now)
	}

	timeout, err := cmd.Flags().GetInt("timeout")
	if err != nil {
		return nil, fmt.Errorf("Invalid value for timeout: %v", err)
	}

	threads, err := rootCmd.Flags().GetInt("threads")
	if err != nil {
		return nil, fmt.Errorf("invalid value for threads: %w", err)
	}

	if threads <= 0 {
		return nil, fmt.Errorf("The number of threads must be positive")
	}

	config := &core.Config{
		HTTP: core.HTTPConfig{
			Insecure: insecure,
			Timeout:  timeout,
		},
		MaxSeverity:    maxSeverity,
		ExportFormats:  exportFormats,
		Urls:           urls,
		ExportFilename: exportFilename,
		SeverityFilter: severityFilter,
		PluginFilter:   pluginFilters,
		Threads:        threads,
	}

	return config, nil
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
