package cmd

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {

}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chopchop",
	Short: "tool for dynamic application security testing on web applications",
	Long: `
 ________                 _________ .__                  _________ .__                    ._.
 /  _____/  ____           \_   ___ \|  |__   ____ ______ \_   ___ \|  |__   ____ ______   | |
/   \  ___ /  _ \   ______ /    \  \/|  |  \ /  _ \\____ \/    \  \/|  |  \ /  _ \\____ \  | |
\    \_\  (  <_> ) /_____/ \     \___|   Y  (  <_> )  |_> >     \___|   Y  (  <_> )  |_> >  \|
 \______  /\____/           \______  /___|  /\____/|   __/ \______  /___|  /\____/|   __/   __
        \/                         \/     \/       |__|           \/     \/       |__|      \/
Link: https://github.com/michelin/ChopChop`,
	SilenceUsage: true,
}

var v string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := setupLogs(os.Stdout, v); err != nil {
			return err
		}
		return nil
	}

	rootCmd.PersistentFlags().StringVarP(&v, "verbosity", "v", logrus.WarnLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(sigs)
		cancel()
	}()
	go func() {
		select {
		case <-sigs:
			log.Warn("\n[!] Keyboard interrupt detected, cleaning up before terminating.")
			cancel()
		case <-ctx.Done():
		}
	}()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Warn(err)
		os.Exit(1)
	}
}

func setupLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	return nil
}
