package pkg

import (
	"fmt"
	"gochopchop/data"
	"log"
	"os"
)

// FormatOutputCSV is a simple wrapper for CSV formatting
func FormatOutputCSV(date string, out []data.Output) {
	f, err := os.OpenFile("./gochopchop_"+date+".csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.WriteString("Domain,endpoint,severity,pluginName,remediation\n")
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < len(out); i++ {
		_, err = f.Write([]byte(out[i].Domain + "," + out[i].TestedURL + "," + out[i].Severity + "," + out[i].PluginName + "," + out[i].Remediation + "\n"))
		if err != nil {
			log.Println(err)
		}
	}

	fmt.Println("Output as csv :" + "./gochopchop_" + date + ".csv")
	f.Close()
}
