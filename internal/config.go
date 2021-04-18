package internal

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

// Config wraps all the parameters to configure
// a scan.
type Config struct {
	HTTP           HTTPConfig
	MaxSeverity    Severity
	SeverityFilter Severity
	ExportFormats  []string
	PluginFilter   []string
	Urls           []string
	ExportFilename string
	Goroutines     int64
}

// HTTPConfig wraps HTTP configurations parameters.
type HTTPConfig struct {
	Insecure bool
	Timeout  int64
}

// ErrNoURL is an error meaning no url is provided.
type ErrNoURL struct{}

func (e ErrNoURL) Error() string {
	return "no urls provided neither with the flag --url-list, -u neither in args"
}

// ErrBothURLAndURLList is an error meaning both args
// and url-file are provided.
type ErrBothURLAndURLList struct{}

func (e ErrBothURLAndURLList) Error() string {
	return "urls provided either with the flag --url-list, -u either in args"
}

// ErrInvalidURLs is an error meaning url(s) is(are)
// invalid.
type ErrInvalidURLs struct {
	urls []string
}

func (e ErrInvalidURLs) Error() string {
	s := "invalid urls: "
	l := len(e.urls)
	for i := 0; i < l; i++ {
		s += e.urls[i]
		if i != l-1 {
			s += ", "
		}
	}
	return s
}

// ErrInvalidExport is an error meaning export(s) is(are)
// not matching allowed exports.
type ErrInvalidExport struct {
	exports []string
}

func (e ErrInvalidExport) Error() string {
	s := "invalid exports: "
	l := len(e.exports)
	for i := 0; i < l; i++ {
		s += e.exports[i]
		if i != l-1 {
			s += ", "
		}
	}
	return s
}

type ErrNegativeField struct {
	Field string
}

func (e ErrNegativeField) Error() string {
	return e.Field + " cannot be negative"
}

// BuildConfig builds the core.Config from provided values.
// Those are supposed to come from the "scan" command flags.
func BuildConfig(insecure bool, export, pluginFilters []string, exportFilename, maxSeverity, severityFilter, urlFile string, threads, timeout int64, args cli.Args) (*Config, error) {
	nArg := args.Len()

	// Check url conditions
	var urls []string
	var invalidUrls []string
	if urlFile == "" {
		if nArg == 0 {
			// There are no args (urls) to chopchop
			return nil, &ErrNoURL{}
		}

		// Check URLs validity
		for i := 0; i < nArg; i++ {
			arg := args.Get(i)
			if !isValidUrl(arg) {
				invalidUrls = append(invalidUrls, arg)
				continue
			}
			urls = append(urls, arg)
		}
	} else {
		// Check there are not args (urls) and an url-file to chopchop
		if nArg != 0 {
			return nil, &ErrBothURLAndURLList{}
		}

		// Open the file containing urls
		c, err := os.Open(urlFile)
		if err != nil {
			return nil, err
		}

		// Read content and add if is a valid url
		scanner := bufio.NewScanner(c)
		for scanner.Scan() {
			url := scanner.Text()
			if !isValidUrl(url) {
				invalidUrls = append(invalidUrls, url)
				continue
			}
			urls = append(urls, url)
		}

		if err := c.Close(); err != nil {
			return nil, err
		}

		// Ensure there were no issues while scanning
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}
	if len(invalidUrls) != 0 {
		return nil, &ErrInvalidURLs{invalidUrls}
	}

	// Check severities
	sevFilter, err := StringToSeverity(severityFilter)
	if err != nil {
		return nil, err
	}
	maxSev, err := StringToSeverity(maxSeverity)
	if err != nil {
		return nil, err
	}

	// Check export
	var invalidExport []string
	for _, e := range export {
		if _, ok := exportersMap[e]; !ok {
			invalidExport = append(invalidExport, e)
		}
	}
	if len(invalidExport) != 0 {
		return nil, &ErrInvalidExport{invalidExport}
	}

	// Check export-filename
	if exportFilename == "" {
		now := time.Now().Format("2006-01-02_15-04-05")
		exportFilename = fmt.Sprintf("gochopchop_%s", now)
	}

	// Check timeout
	if timeout < 0 {
		return nil, &ErrNegativeField{"timeout"}
	}

	// Check threads
	if threads < 0 {
		return nil, &ErrNegativeField{"threads"}
	}

	// Check insecure       => always fine
	// Check plugin-filters => always fine ?

	// Build config
	config := &Config{
		HTTP: HTTPConfig{
			Insecure: insecure,
			Timeout:  timeout,
		},
		MaxSeverity:    maxSev,
		SeverityFilter: sevFilter,
		ExportFormats:  export,
		PluginFilter:   pluginFilters,
		Urls:           urls,
		ExportFilename: exportFilename,
		Goroutines:     threads,
	}

	return config, nil
}

func isValidUrl(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}
