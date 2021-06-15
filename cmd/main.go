package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"sort"
	"syscall"

	"github.com/michelin/gochopchop/internal"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	cliLogo = `
  ________                 _________ .__                  _________ .__                    ._.
 /  _____/  ____           \_   ___ \|  |__   ____ ______ \_   ___ \|  |__   ____ ______   | |
/   \  ___ /  _ \   ______ /    \  \/|  |  \ /  _ \\____ \/    \  \/|  |  \ /  _ \\____ \  | |
\    \_\  (  <_> ) /_____/ \     \___|   Y  (  <_> )  |_> >     \___|   Y  (  <_> )  |_> >  \|
 \______  /\____/           \______  /___|  /\____/|   __/ \______  /___|  /\____/|   __/   __
        \/                         \/     \/       |__|           \/     \/       |__|      \/
`
	AppHelpTemplate = cliLogo + `
{{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

Usage:
   chopchop [command]{{"\n"}}

{{- if .Description}}
DESCRIPTION:
   {{.Description | nindent 3 | trim}}{{end}}

{{- if .VisibleCommands}}
Available Commands:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{else}}{{range .VisibleCommands}}
   {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

Flags:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}

Use "chopchop [command] --help" for more information about a command.
`
	CommandHelpTemplate = `{{.Usage}}

Usage:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [flags]{{end}}{{end}}{{if .VisibleFlags}}

Flags:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`
)

func flagsMdw(flags []cli.Flag) []cli.Flag {
	// Build shared flags
	f := []cli.Flag{
		&cli.IntFlag{
			Name:  "threads",
			Usage: "number of threads (goroutines to be exact)",
			Value: 1,
		},
		&cli.StringFlag{
			Name:    "verbosity",
			Aliases: []string{"v"},
			Usage:   "log level (debug, info, warn, error, fatal, panic)",
			Value:   "warning",
		},
	}

	return append(flags, f...)
}

func cliMdw(f func(*cli.Context) error) func(*cli.Context) error {
	return func(c *cli.Context) error {
		// Setup logs
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		lvl, err := logrus.ParseLevel(c.String("verbosity"))
		if err != nil {
			return err
		}
		logrus.SetLevel(lvl)
		logrus.Debug("verbosity:", lvl)

		// Call the wrapped cli func
		return f(c)
	}
}

func main() {
	app := &cli.App{
		Name:  "ChopChop",
		Usage: "CLI tool to help developers scanning endpoints and identifying exposition of sensitive services/files/folders.\nhttps://github.com/michelin/ChopChop.",
		Commands: []*cli.Command{
			{
				Name:   "plugins",
				Usage:  "list checks of configuration file",
				Action: cliMdw(cmdPlugins),
				Flags: flagsMdw([]cli.Flag{
					&cli.StringFlag{
						Name:    "severity",
						Aliases: []string{"s"},
						Usage:   "severity option for list tag",
						Value:   "Informational",
					},
					&cli.StringFlag{
						Name:    "signatures",
						Aliases: []string{"c"},
						Usage:   "path to signature file",
						Value:   "chopchop.yml",
					},
				}),
			}, {
				Name:   "scan",
				Usage:  "scan endpoints to check if services/files/folders are exposed",
				Action: cliMdw(cmdScan),
				Flags: flagsMdw([]cli.Flag{
					&cli.StringSliceFlag{
						Name:    "export",
						Aliases: []string{"e"},
						Usage:   "export of the output (" + internal.ExportersList() + ")",
						Value:   &cli.StringSlice{},
					},
					&cli.StringFlag{
						Name:  "export-filename",
						Usage: "filename for export files",
						Value: "",
					},
					&cli.BoolFlag{
						Name:    "insecure",
						Aliases: []string{"k"},
						Usage:   "check SSL certificate",
						Value:   false,
					},
					&cli.StringFlag{
						Name:    "max-severity",
						Aliases: []string{"b"},
						Usage:   "block the CI pipeline if severity is over or equal specified flag",
						Value:   "Informational",
					},
					&cli.StringSliceFlag{
						Name:  "plugin-filters",
						Usage: "filter by the name of the plugin (engine will only check for plugin with the same name)",
						Value: &cli.StringSlice{},
					},
					&cli.StringFlag{
						Name:  "severity-filter",
						Usage: "filter by severity (engine will check for same severity checks)",
						Value: "Informational",
					},
					&cli.StringFlag{
						Name:    "signatures",
						Aliases: []string{"c"},
						Usage:   "path to signature file",
						Value:   "chopchop.yml",
					},
					&cli.IntFlag{
						Name:    "timeout",
						Aliases: []string{"t"},
						Usage:   "timeout (in s) for the HTTP requests",
						Value:   10,
					},
					&cli.StringFlag{
						Name:    "url-file",
						Aliases: []string{"u"},
						Usage:   "path to a specified file containing urls to test",
						Value:   "",
					},
				}),
			},
		},
	}
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate

	// Setup stop signals
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(sigs)
		cancel() // Triggers the <-ctx.Done() in the following goroutine
	}()
	go func() {
		select {
		case <-sigs:
			logrus.Warn("Keyboard interrupt detected.")
			cancel()
			os.Exit(1)
		case <-ctx.Done():
			return
		}
	}()

	err := app.RunContext(ctx, os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

func cmdScan(c *cli.Context) error {
	// Build the config
	insecure := c.Bool("insecure")
	exprt := c.StringSlice("export")
	pluginFilters := c.StringSlice("plugin-filters")
	exportFilename := c.String("export-filename")
	maxSeverity := c.String("max-severity")
	severityFilter := c.String("severity-filter")
	urlFile := c.String("url-file")
	timeout := c.Int("timeout")
	threads := c.Int("threads")
	args := c.Args()

	logrus.Debug("insecure:", insecure)
	logrus.Debug("export:", exprt)
	logrus.Debug("plugin-filters:", pluginFilters)
	logrus.Debug("export-filename:", exportFilename)
	logrus.Debug("max-severity:", maxSeverity)
	logrus.Debug("severity-filter:", severityFilter)
	logrus.Debug("url-file:", urlFile)
	logrus.Debug("timeout:", timeout)
	logrus.Debug("threads:", threads)
	logrus.Debug("args:", args)

	var urlFileReader io.Reader
	if urlFile != "" {
		var err error
		urlFileReader, err = os.Open(urlFile)
		if err != nil {
			return err
		}
	}

	config, err := internal.BuildConfig(insecure, exprt, pluginFilters, exportFilename, maxSeverity, severityFilter, urlFileReader, threads, timeout, args.Slice())
	if err != nil {
		return err
	}

	// Parse signatures
	signatures := c.String("signatures")

	signFile, err := internal.ReaderFromFile(signatures)
	if err != nil {
		return err
	}
	sign, err := internal.ParseSignatures(signFile)
	if err != nil {
		return err
	}

	// Build the CoreScanner
	scanner, err := internal.NewCoreScanner(config, sign)
	if err != nil {
		return err
	}

	// Start the scan
	results, dur, err := internal.Scan(scanner, config.Urls, c.Done())
	if err != nil {
		return err
	}
	logrus.Info("Scan execution time: ", dur)

	// Sort and export the results
	sort.Stable(results)
	err = internal.ExportResults(results, config, exportFilename)
	if err != nil {
		return err
	}

	return nil
}

func cmdPlugins(c *cli.Context) error {
	// Parse signatures
	signatures := c.String("signatures")

	logrus.Debug("signatures:", signatures)

	signFile, err := internal.ReaderFromFile(signatures)
	if err != nil {
		return err
	}
	sign, err := internal.ParseSignatures(signFile)
	if err != nil {
		return err
	}

	// Parse severity
	severity := c.String("severity")

	sev, err := internal.StringToSeverity(severity)
	if err != nil {
		return err
	}
	sevStr, _ := sev.String()

	// Print signatures in stdout
	internal.PrintSignatures(sign, sevStr, os.Stdout)

	return nil
}
