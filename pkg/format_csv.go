package pkg

import (
	"fmt"
	"gochopchop/data"
	"log"
	"os"
)

// WriteCSVOutput is a simple wrapper for CSV formatting
func WriteCSVOutput(fileResults string, out []data.Output) {
	f, err := os.OpenFile(fileResults, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.WriteString("Domain,uri,severity,pluginName,remediation\n")
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < len(out); i++ {
		_, err = f.Write([]byte(out[i].Domain + "," + out[i].TestedURL + "," + out[i].Severity + "," + out[i].PluginName + "," + out[i].Remediation + "\n"))
		if err != nil {
			log.Println(err)
		}
	}

	fmt.Println("Successfuly written output in: " + fileResults)
	f.Close()
}
