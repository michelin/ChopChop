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
	rootCmd.AddCommand(pluginCmd)
	pluginCmd.Flags().StringP("config-file", "c", "chopchop.yml", "path to config/data file") // --config-file ou -c
	pluginCmd.Flags().StringP("severity", "s", "", "severity option for list tag")            // --severity ou -s
}

var pluginCmd = &cobra.Command{
	Use:   "plugins",
	Short: "list checks of configuration file",
	Args:  pluginCheckArgsAndFlags,
	Run:   app.List,
}

func pluginCheckArgsAndFlags(cmd *cobra.Command, args []string) error {
	//validate config filepath
	configFile, err := cmd.Flags().GetString("config-file")
	if err != nil {
		return fmt.Errorf("invalid value for configFile: %v", err)
	}
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Fatal("Filepath of config file is not valid")
	}
	if !strings.HasSuffix(configFile, ".yml") {
		log.Fatal("Config file needs to be a yaml file")
	}
	//Validate severity input
	severity, err := cmd.Flags().GetString("severity")
	if err != nil {
		return fmt.Errorf("invalid value for severity: %v", err)
	}
	if severity != "" {
		if severity == "High" || severity == "Medium" || severity == "Low" || severity == "Informational" {
			fmt.Println("Display only check of severity : " + severity)
		} else {
			log.Fatal(" ------ Unknown severity type : " + severity + " . Only Informational / Low / Medium / High are valid severity types.")
		}
	}
	if err := cmd.Flags().Set("config-file", configFile); err != nil {
		return fmt.Errorf("error while setting filepath of config file")
	}
	return nil
}
