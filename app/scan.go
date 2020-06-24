package app

import (
	"bufio"
	"fmt"
	"gochopchop/data"
	"gochopchop/pkg"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// SeverityType is basically an enum and values can be from Info, Low, Medium and High
type SeverityType string

const (
	// Informational will be the default severityType
	Informational SeverityType = "Informational"
	// Low severity
	Low = "Low"
	// Medium severity
	Medium = "Medium"
	// High severity (highest rating)
	High = "High"
)

// Config struct to load the configuration from the YAML file
type Config struct {
	Insecure bool `yaml:"insecure"`
	Plugins  []struct {
		URI    string `yaml:"uri"`
		Checks []struct {
			Match       []*string     `yaml:"match"`
			AllMatch    []*string     `yaml:"all_match"`
			StatusCode  *int32        `yaml:"status_code"`
			PluginName  string        `yaml:"name"`
			Remediation *string       `yaml:"remediation"`
			Severity    *SeverityType `yaml:"severity"`
			Description *string       `yaml:"description"`
			NoMatch     []*string     `yaml:"no_match"`
			Headers     []*string     `yaml:"headers"`
		} `yaml:"checks"`
	} `yaml:"plugins"`
}

// Scan of domain via url
func Scan(cmd *cobra.Command, args []string) {
	url, _ := cmd.Flags().GetString("url")
	insecure, _ := cmd.Flags().GetBool("insecure")
	csv, _ := cmd.Flags().GetBool("csv")
	json, _ := cmd.Flags().GetBool("json")
	urlFile, _ := cmd.Flags().GetString("url-file")
	configFile, _ := cmd.Flags().GetString("config-file")
	suffix, _ := cmd.Flags().GetString("suffix")
	prefix, _ := cmd.Flags().GetString("prefix")
	blockedFlag, _ := cmd.Flags().GetString("block")

	var tmpURL string
	var urlList []string

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

	y := Config{}
	if err = yaml.Unmarshal([]byte(dataCfg), &y); err != nil {
		log.Fatal(err)
	}
	// If flag insecure isn't specified, check yaml file if it's specified in it
	if insecure {
		fmt.Println("Launching scan without validating the SSL certificate")
	} else {
		insecure = y.Insecure
	}

	CheckStructFields(y)
	hit := false
	block := false
	currentTime := time.Now()
	date := currentTime.Format("2006-01-02_15-04-05")
	out := []data.Output{}

	for i := 0; i < len(urlList); i++ {
		fmt.Print("Testing domain : ")
		fmt.Println(prefix + urlList[i] + suffix)
		for index, plugin := range y.Plugins {
			_ = index
			tmpURL = prefix + urlList[i] + suffix + fmt.Sprint(plugin.URI)
			httpResponse, err := pkg.HTTPGet(insecure, tmpURL)
			if err != nil {
				_ = errors.Wrap(err, "Timeout of HTTP Request")
			}

			if httpResponse != nil {
				for index, check := range plugin.Checks {
					_ = index
					answer := pkg.ResponseAnalysis(httpResponse, check.StatusCode, check.Match, check.AllMatch, check.NoMatch, check.Headers)
					if answer {
						hit = true
						if BlockCI(blockedFlag, *check.Severity) {
							block = true
						}
						out = append(out, data.Output{
							Domain:      urlList[i],
							PluginName:  check.PluginName,
							TestedURL:   plugin.URI,
							Severity:    string(*check.Severity),
							Remediation: *check.Remediation,
						})
					}
				}
			} else {
				fmt.Println("Server refused the connection for URL : " + tmpURL)
				continue
			}
			_ = httpResponse.Body.Close()
		}
	}
	if hit {
		pkg.FormatOutputTable(out)
		if json {
			outputJSON := pkg.AddVulnToOutputJSON(out)
			pkg.CreateFileJSON(date, outputJSON)
		}
		if csv {
			pkg.FormatOutputCSV(date, out)
		}
		if blockedFlag != "" {
			if block {
				os.Exit(1)
			} else {
				fmt.Println("No critical vulnerabilities found...")
				os.Exit(0)
			}
		}
		os.Exit(1)
	} else {
		fmt.Println("No vulnerabilities found. Exiting...")
		os.Exit(0)
	}
}

// BlockCI function will allow the user to return a different status code depending on the highest severity that has triggered
func BlockCI(severity string, severityType SeverityType) bool {
	switch severity {
	case "High":
		if severityType == High {
			return true
		}
	case "Medium":
		if severityType == High || severityType == Medium {
			return true
		}
	case "Low":
		if severityType == High || severityType == Medium || severityType == Low {
			return true
		}
	case "Informational":
		if severityType == High || severityType == Medium || severityType == Low || severityType == Informational {
			return true
		}
	}
	return false
}

// CheckStructFields will parse the YAML configuration file
func CheckStructFields(conf Config) {
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
				if err := SeverityType.IsValid(*check.Severity); err != nil {
					log.Fatal(" ------ Unknown severity type : " + string(*check.Severity) + " . Only Informational / Low / Medium / High are valid severity types.")
				}
			}
		}
	}
}

// IsValid will verify that the severityType is part of the enum previously declared
func (st SeverityType) IsValid() error {
	switch st {
	case Informational, Low, Medium, High:
		return nil
	}
	return errors.New("Invalid Severity type. Please Check yaml config file")
}
