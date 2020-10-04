package pkg

import (
	"encoding/json"
	"fmt"
	"gochopchop/data"
	"log"
	"os"
)

// OutputJSON struct
type OutputJSON struct {
	TestedDomains []TestedDomains `json:"domains"`
}

// TestedDomains struct
type TestedDomains struct {
	TestedDomain string       `json:"domain"`
	TestedUrls   []TestedURLs `json:"urls"`
}

// TestedURLs struct
type TestedURLs struct {
	TestedURL   string `json:"url,omitempty"`
	PluginName  string `json:"plugin_name,omitempty"`
	Severity    string `json:"severity,omitempty"`
	Remediation string `json:"remediation,omitempty"`
}

// AddVulnToOutputJSON will add the vuln to struct output
func AddVulnToOutputJSON(out []data.Output) OutputJSON {
	jsonOut := OutputJSON{}

	for i := 0; i < len(out); i++ {
		added := false
		// Check if domain already exist - if yes append infos
		for y := 0; y < len(jsonOut.TestedDomains); y++ {
			if jsonOut.TestedDomains[y].TestedDomain == out[i].Domain {
				jsonOut.TestedDomains[y].TestedUrls = append(jsonOut.TestedDomains[y].TestedUrls, TestedURLs{
					TestedURL:   out[i].TestedURL,
					PluginName:  out[i].PluginName,
					Severity:    out[i].Severity,
					Remediation: out[i].Remediation,
				})
				added = true
			}
		}
		if !added {
			// If domain not found, create it
			jsonOut.TestedDomains = append(jsonOut.TestedDomains, TestedDomains{
				TestedDomain: out[i].Domain,
				TestedUrls:   nil,
			})
			jsonOut.TestedDomains[len(jsonOut.TestedDomains)-1].TestedUrls = append(jsonOut.TestedDomains[len(jsonOut.TestedDomains)-1].TestedUrls, TestedURLs{
				TestedURL:   out[i].TestedURL,
				PluginName:  out[i].PluginName,
				Severity:    out[i].Severity,
				Remediation: out[i].Remediation,
			})
		}
	}
	return jsonOut
}

// WriteJSONOutput will save the output to a JSON file
func WriteJSONOutput(fileResults string, out OutputJSON) {
	f, err := os.OpenFile(fileResults, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	file, _ := json.MarshalIndent(&out, "", " ")

	_, err = f.Write([]byte(file))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfuly written output in:" + fileResults)
	f.Close()
}
