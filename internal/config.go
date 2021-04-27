package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"
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
	Goroutines     int
}

// HTTPConfig wraps HTTP configurations parameters.
type HTTPConfig struct {
	Insecure bool
	Timeout  int
}

// ErrNoURL is an error meaning no url is provided.
var ErrNoURL = errors.New("no urls provided neither with the flag --url-list, -u neither in args")

// ErrBothURLAndURLList is an error meaning both args
// and url-file are provided.
var ErrBothURLAndURLList = errors.New("urls provided either with the flag --url-list, -u either in args")

// ErrInvalidURLs is an error meaning url(s) is(are)
// invalid.
type ErrInvalidURLs struct {
	URLs []string
}

func (e ErrInvalidURLs) Error() string {
	s := "invalid URLs: "
	l := len(e.URLs)
	for i := 0; i < l; i++ {
		s += e.URLs[i]
		if i != l-1 {
			s += ", "
		}
	}
	return s
}

// ErrInvalidExport is an error meaning export(s) is(are)
// not matching allowed exports.
type ErrInvalidExport struct {
	Exports []string
}

func (e ErrInvalidExport) Error() string {
	s := "invalid exports: "
	l := len(e.Exports)
	for i := 0; i < l; i++ {
		s += e.Exports[i]
		if i != l-1 {
			s += ", "
		}
	}
	return s
}

// ErrFailedOperationOnField is an error meaning
// a field has failed to pass an operation on a value.
type ErrFailedOperationOnField struct {
	Field     string
	Operation string
	Value     int
}

func (e ErrFailedOperationOnField) Error() string {
	return e.Field + " failed to be " + e.Operation + " (specified " + strconv.Itoa(e.Value) + ")"
}

// BuildConfig builds the core.Config from provided values.
// Those are supposed to come from the "scan" command flags.
func BuildConfig(insecure bool, export, pluginFilters []string, exportFilename, maxSeverity, severityFilter string, urlFile io.Reader, threads, timeout int, args []string) (*Config, error) {
	nArg := len(args)

	// Check insecure       => always fine

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

	// Check severities
	sevFilter, err := StringToSeverity(severityFilter)
	if err != nil {
		return nil, err
	}
	maxSev, err := StringToSeverity(maxSeverity)
	if err != nil {
		return nil, err
	}

	// Check url conditions
	var urls []string
	var invalidUrls []string
	if urlFile == nil {
		if nArg == 0 {
			// There are no args (urls) to chopchop
			return nil, ErrNoURL
		}

		// Check URLs validity
		for i := 0; i < nArg; i++ {
			arg := args[i]
			if !isValidUrl(arg) {
				invalidUrls = append(invalidUrls, arg)
				continue
			}
			urls = append(urls, arg)
		}
	} else {
		// Check there are not args (urls) and an url-file to chopchop
		if nArg != 0 {
			return nil, ErrBothURLAndURLList
		}

		// Read content and add if is a valid url
		scanner := bufio.NewScanner(urlFile)
		for scanner.Scan() {
			url := scanner.Text()
			if !isValidUrl(url) {
				invalidUrls = append(invalidUrls, url)
				continue
			}
			urls = append(urls, url)
		}

		// Ensure there were no issues while scanning
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}
	if len(invalidUrls) != 0 {
		return nil, &ErrInvalidURLs{invalidUrls}
	}

	// Check threads
	if threads <= 0 {
		return nil, &ErrFailedOperationOnField{"threads", "<=0", threads}
	}

	// Check timeout
	if timeout < 0 {
		return nil, &ErrFailedOperationOnField{"timeout", "<0", timeout}
	}

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
