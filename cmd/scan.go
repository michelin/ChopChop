package cmd

import (
	"fmt"
	"gochopchop/app"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringP("url", "u", "", "url to scan")                                                                                                     // --url OU -u
	scanCmd.Flags().StringP("config-file", "c", "chopchop.yml", "path to config/data file")                                                                    // --config-file ou -c
	scanCmd.Flags().BoolP("insecure", "i", false, "Check SSL certificate")                                                                                     // --insecure ou -n
	scanCmd.Flags().StringP("url-file", "f", "", "path to a specified file containing urls to test")                                                           // --uri-file ou -f
	scanCmd.Flags().StringP("suffix", "s", "", "Add suffix to urls when flag url-file is specified")                                                           // --suffix ou -s
	scanCmd.Flags().StringP("prefix", "p", "", "Add prefix to urls when flag url-file is specified")                                                           // --prefix ou -p
	scanCmd.Flags().Int32P("timeout", "t", 10, "Timeout for the HTTP requests (default: 10s)")                                                                 // --timeout ou -t
	scanCmd.Flags().StringP("block", "b", "", "Block pipeline if severity is over or equal specified flag")                                                    // --block ou -b
	scanCmd.Flags().StringP("signature-name", "", "", "Filter by signature names (engine will check if words are contained), can use comma for multiple ones") // --signature-name
	scanCmd.Flags().StringP("severity", "", "", "Filter by severity (engine will check for same severity checks)")                                             // --severity
	scanCmd.Flags().BoolP("csv", "", false, "output as a csv file")                                                                                            //--csv
	scanCmd.Flags().BoolP("json", "", false, "output as a json file")                                                                                          //--json
	scanCmd.Flags().StringP("csv-file", "", "results.csv", "output as a csv file (Default: results.csv)")                                                      //--csv mydomain.csv
	scanCmd.Flags().StringP("json-file", "", "results.json", "output as a json file (Default: results.json)")                                                  //--json mydomain.json
	scanCmd.Flags().BoolP("verbose", "v", false, "Verbose mode")                                                                                               //--verbose ou -v
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan URL endpoints to check if services/files/folders are exposed to the Internet",
	Args:  scanCheckArgsAndFlags,
	Run:   app.Scan,
}

func scanCheckArgsAndFlags(cmd *cobra.Command, args []string) error {
	// validate url
	url, err := cmd.Flags().GetString("url")
	if err != nil {
		return fmt.Errorf("invalid value for url: %v", err)
	}
	configFile, err := cmd.Flags().GetString("config-file")
	if err != nil {
		return fmt.Errorf("invalid value for configFile: %v", err)
	}
	urlFile, err := cmd.Flags().GetString("url-file")
	if err != nil {
		return fmt.Errorf("invalid value for urlFile: %v", err)
	}
	suffix, err := cmd.Flags().GetString("suffix")
	if err != nil {
		return fmt.Errorf("invalid value for suffix: %v", err)
	}
	prefix, err := cmd.Flags().GetString("prefix")
	if err != nil {
		return fmt.Errorf("invalid value for prefix: %v", err)
	}
	block, err := cmd.Flags().GetString("block")
	if err != nil {
		return fmt.Errorf("invalid value for block: %v", err)
	}
	// if url != "" {
	// 	if !strings.HasPrefix(url, "http") {
	// 		// If http or https not specified, return fatal log
	// 		return fmt.Errorf("URL needs a specified prefix :  http:// or https://")
	// 	}
	// }
	if suffix != "" || prefix != "" {
		if urlFile == "" {
			return fmt.Errorf("suffix or prefix flags can't be assigned if flag url-file is not specified")
		}
	}
	if block != "" {
		if block == "High" || block == "Medium" || block == "Low" || block == "Informational" {
			fmt.Println("Block pipeline if severity is over or equal : " + block)
		} else {
			log.Fatal(" ------ Unknown severity type : " + block + " . Only Informational / Low / Medium / High are valid severity types.")
		}
	}
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("filepath of config file is not valid")
	}
	if !strings.HasSuffix(configFile, ".yml") {
		return fmt.Errorf("config file needs to be a yaml file")
	}
	if urlFile != "" {
		if _, err := os.Stat(urlFile); os.IsNotExist(err) {
			return fmt.Errorf("filepath of url file is not valid")
		}
	}
	if err := cmd.Flags().Set("config-file", configFile); err != nil {
		return fmt.Errorf("error while setting filepath of config file")
	}
	if err := cmd.Flags().Set("url-file", urlFile); err != nil {
		return fmt.Errorf("error while setting filepath of url file")
	}
	if err := cmd.Flags().Set("url", url); err != nil {
		return fmt.Errorf("error while setting url")
	}
	if err := cmd.Flags().Set("suffix", suffix); err != nil {
		return fmt.Errorf("error while setting suffix")
	}
	if err := cmd.Flags().Set("prefix", prefix); err != nil {
		return fmt.Errorf("error while setting prefix")
	}
	if err := cmd.Flags().Set("block", block); err != nil {
		return fmt.Errorf("error while setting block flag")
	}
	return nil
}
