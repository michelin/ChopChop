package app

import (
	"bufio"
	"fmt"
	"gochopchop/data"
	"gochopchop/pkg"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Verbose Verbose function
func Verbose(message string, verbose bool) {
	if verbose {
		fmt.Println("[verbose] " + message)
	}
}

// Scan of domain via url
func Scan(cmd *cobra.Command, args []string) {
	start := time.Now()

	url, _ := cmd.Flags().GetString("url")
	insecure, _ := cmd.Flags().GetBool("insecure")
	csv, _ := cmd.Flags().GetBool("csv")
	json, _ := cmd.Flags().GetBool("json")
	csvFile, _ := cmd.Flags().GetString("csv-file")
	jsonFile, _ := cmd.Flags().GetString("json-file")
	signatureName, _ := cmd.Flags().GetString("signature-name")
	severity, _ := cmd.Flags().GetString("severity")
	urlFile, _ := cmd.Flags().GetString("url-file")
	configFile, _ := cmd.Flags().GetString("config-file")
	suffix, _ := cmd.Flags().GetString("suffix")
	prefix, _ := cmd.Flags().GetString("prefix")
	httpRequestTimeout, _ := cmd.Flags().GetInt("timeout")
	blockedFlag, _ := cmd.Flags().GetString("block")
	verbose, _ := cmd.Flags().GetBool("verbose")

	var tmpURL string
	var urlList []string

	if url == "" && urlFile == "" {
		log.Fatal("`url` or either `url-file` have been specified! Use `scan -help` for usage")
	}

	cfg, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}

	defer cfg.Close()
	dataCfg, err := ioutil.ReadAll(cfg)

	if url != "" {
		urlList = append(urlList, url)
	}

	if urlFile != "" {
		urlFileContent, err := os.Open(urlFile)
		if err != nil {
			log.Fatal(err)
		}
		defer urlFileContent.Close()

		scanner := bufio.NewScanner(urlFileContent)
		for scanner.Scan() {
			urlList = append(urlList, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	y := data.Config{}
	if err = yaml.Unmarshal([]byte(dataCfg), &y); err != nil {
		log.Fatal(err)
	}
	// If flag insecure isn't specified, check yaml file if it's specified in it
	if insecure {
		Verbose("Launching scan without validating the SSL certificate", verbose)
	} else {
		insecure = y.Insecure
	}

	CheckStructFields(y)
	hit := false
	block := false
	out := []data.Output{}

	var wg sync.WaitGroup
	wg.Add(len(urlList))

	tags := strings.Split(signatureName, ",")

	for i := 0; i < len(urlList); i++ {
		go func(domain string) {
			defer wg.Done()
			Verbose("Testing domain : "+prefix+domain+suffix, verbose)
			for _, plugin := range y.Plugins {

				var uris []string
				// check the present of `uri` and `uris` which has to be forbidden
				if len(plugin.ListOfURI) > 0 && plugin.URI != "" {
					log.Fatal("You can't both have `uri` and `uris` specified for URI: " + plugin.URI)
				}

				if plugin.URI != "" {
					uris = make([]string, 1)
					uris[0] = plugin.URI
				} else {
					uris = plugin.ListOfURI
				}

				for j := 0; j < len(uris); j++ {
					canExecutePluginByName := isPluginAuthorized(plugin, tags)
					if !canExecutePluginByName {
						Verbose("Skipping signature rule with URI: "+uris[j], verbose)
					}

					canExecutePluginBySeverity := isPluginAuthorizedBySeverity(plugin, severity)
					if !canExecutePluginBySeverity {
						Verbose("Skipping signature rule with URI: "+uris[j], verbose)
					}

					if canExecutePluginByName && canExecutePluginBySeverity {
						tmpURL = prefix + domain + suffix + fmt.Sprint(uris[j])
						if plugin.QueryString != "" {
							tmpURL += "?" + plugin.QueryString
						}

						// By default we follow HTTP redirects
						followRedirects := true
						// But for each plugin we can override and don't follow HTTP redirects
						if plugin.FollowRedirects != nil && *plugin.FollowRedirects == false {
							followRedirects = false
						}

						Verbose("Testing URL: "+tmpURL, verbose)
						httpResponse, err := pkg.HTTPGet(insecure, tmpURL, followRedirects, httpRequestTimeout)
						if err != nil {
							fmt.Println(err)
						}

						if httpResponse != nil {
							for _, check := range plugin.Checks {
								if severity != "" && check.Severity.String() != severity {
									continue // break if the severity is defined and different from the one we specified
								}
								answer := pkg.ResponseAnalysis(httpResponse, check)
								if answer {
									Verbose("[!] Hit found!\n\tURL: "+tmpURL+"\n\tPlugin: "+check.PluginName+"\n\tSeverity: "+string(*check.Severity), verbose)
									hit = true
									if BlockCI(blockedFlag, *check.Severity) {
										block = true
									}
									out = append(out, data.Output{
										Domain:      domain,
										PluginName:  check.PluginName,
										TestedURL:   uris[j],
										Severity:    string(*check.Severity),
										Remediation: *check.Remediation,
									})
								}
							}
							_ = httpResponse.Body.Close()
						}
					}
				}
			}
		}(urlList[i])
	}

	wg.Wait()

	if hit {
		pkg.FormatOutputTable(out)
		if json {
			outputJSON := pkg.AddVulnToOutputJSON(out)
			pkg.WriteJSONOutput(jsonFile, outputJSON)
		}
		if csv {
			pkg.WriteCSVOutput(csvFile, out)
		}
	}

	elapsed := time.Since(start)
	Verbose(fmt.Sprintf("Scan execution time: %s", elapsed), verbose)

	// return EXIT_SUCCESS if
	// 1. no hit
	// 2. no vulns >= the cricity we're looking for
	if hit {
		if blockedFlag != "" && !block {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	} else {
		fmt.Println("No vulnerabilities found.")
	}
}

// isPluginAuthorizedBySeverity returns `true` if there's at least one check with the same severity
func isPluginAuthorizedBySeverity(signature data.Signature, severity string) bool {
	// if the flag is not properly set by the user, return true and execute them all
	if severity == "" {
		return true
	}
	for _, check := range signature.Checks {
		if check.Severity.String() == severity {
			return true
		}
	}
	return false
}

// isPluginAuthorized returns `true` if there's at least one tag in a check name
func isPluginAuthorized(signature data.Signature, tags []string) bool {
	// if the flag is not properly set by the user, return true and execute them all
	if tags[0] == "" {
		return true
	}
	for j := 0; j < len(tags); j++ {
		for _, check := range signature.Checks {
			// fmt.Println(check.PluginName)
			// fmt.Println(tags[j])

			// if there's one tag in one of the `checks`, we do the request
			// otherwise, we don't
			if strings.Contains(check.PluginName, tags[j]) {
				return true
			}
		}
	}
	return false
}

// BlockCI function will allow the user to return a different status code depending on the highest severity that has triggered
func BlockCI(severity string, severityType data.SeverityType) bool {
	switch severity {
	case "High":
		if severityType == data.High {
			return true
		}
	case "Medium":
		if severityType == data.High || severityType == data.Medium {
			return true
		}
	case "Low":
		if severityType == data.High || severityType == data.Medium || severityType == data.Low {
			return true
		}
	case "Informational":
		if severityType == data.High || severityType == data.Medium || severityType == data.Low || severityType == data.Informational {
			return true
		}
	}
	return false
}

// CheckStructFields will parse the YAML configuration file
func CheckStructFields(conf data.Config) {
	for index, plugin := range conf.Plugins {
		_ = index
		for index, check := range plugin.Checks {
			_ = index
			if check.Description == nil {
				log.Fatal("Missing description field in " + check.PluginName + " plugin checks. Stopping execution.")
			}
			if check.Remediation == nil {
				log.Fatal("Missing remediation field in " + check.PluginName + " plugin checks. Stopping execution.")
			}
			if check.Severity == nil {
				log.Fatal("Missing severity field in " + check.PluginName + " plugin checks. Stopping execution.")
			} else {
				if err := data.SeverityType.IsValid(*check.Severity); err != nil {
					log.Fatal(" ------ Unknown severity type : " + string(*check.Severity) + " . Only Informational / Low / Medium / High are valid severity types.")
				}
			}
		}
	}
}
