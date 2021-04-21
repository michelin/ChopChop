package internal_test

import (
	"io"
	"reflect"
	"testing"

	"github.com/michelin/gochopchop/internal"
)

type FakeReadCloser struct {
	data      []byte
	readIndex int64
}

func NewFakeReadCloser(toRead string) *FakeReadCloser {
	return &FakeReadCloser{data: []byte(toRead)}
}

func (f *FakeReadCloser) Read(p []byte) (n int, err error) {
	if f.readIndex >= int64(len(f.data)) {
		err = io.EOF
		return
	}

	n = copy(p, f.data[f.readIndex:])
	f.readIndex += int64(n)
	return
}

func (f *FakeReadCloser) Close() error {
	return nil
}

var _ = (io.ReadCloser)(&FakeReadCloser{})

func TestBuildConfig(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		Insecure       bool
		Export         []string
		PluginFilters  []string
		ExportFilename string
		MaxSeverity    string
		SeverityFilter string
		URLFile        io.Reader
		Threads        int64
		Timeout        int64
		Args           []string
		ExpectedConfig *internal.Config
		ExpectedErr    error
	}{
		"inexisting-exporter": {
			Export:         []string{"inexisting-exporter", "other-inexisting-exporter"},
			ExpectedConfig: nil,
			ExpectedErr:    &internal.ErrInvalidExport{[]string{"inexisting-exporter", "other-inexisting-exporter"}},
		},
		"invalid-severityfilter": {
			SeverityFilter: "invalid-severity",
			ExpectedConfig: nil,
			ExpectedErr:    &internal.ErrInvalidSeverity{"invalid-severity"},
		},
		"invalid-maxseverity": {
			SeverityFilter: "Low",
			MaxSeverity:    "invalid-severity",
			ExpectedConfig: nil,
			ExpectedErr:    &internal.ErrInvalidSeverity{"invalid-severity"},
		},
		"nil-urlfile-no-args": {
			SeverityFilter: "Low",
			MaxSeverity:    "Low",
			URLFile:        nil,
			Args:           []string{},
			ExpectedConfig: nil,
			ExpectedErr:    internal.ErrNoURL,
		},
		"nil-urlfile-args-invalid-urls": {
			SeverityFilter: "Low",
			MaxSeverity:    "Low",
			URLFile:        nil,
			Args:           []string{"https://www.michelin.com/", "gochopchop", "ChopChop"},
			ExpectedConfig: nil,
			ExpectedErr:    &internal.ErrInvalidURLs{[]string{"gochopchop", "ChopChop"}},
		},
		"urlfile-and-args": {
			SeverityFilter: "Low",
			MaxSeverity:    "Low",
			URLFile:        NewFakeReadCloser(""),
			Args:           []string{""},
			ExpectedConfig: nil,
			ExpectedErr:    internal.ErrBothURLAndURLList,
		},
		"urlfile-no-args-invalid-urls": {
			SeverityFilter: "Low",
			MaxSeverity:    "Low",
			URLFile:        NewFakeReadCloser("https://www.michelin.com/\ngochopchop\nChopChop\n"),
			Args:           []string{},
			ExpectedConfig: nil,
			ExpectedErr:    &internal.ErrInvalidURLs{[]string{"gochopchop", "ChopChop"}},
		},
		"url-file-args-invalid-urls-scanner-err": {
			SeverityFilter: "Low",
			MaxSeverity:    "Low",
			URLFile:        &FailingReadCloser{},
			Args:           []string{},
			ExpectedConfig: nil,
			ExpectedErr:    errFake,
		},
		"negative-threads": {
			SeverityFilter: "Low",
			MaxSeverity:    "Low",
			Threads:        -1,
			Args:           []string{"https://www.michelin.com/"},
			ExpectedConfig: nil,
			ExpectedErr:    &internal.ErrFailedOperationOnField{"threads", "<=0", -1},
		},
		"zero-threads": {
			SeverityFilter: "Low",
			MaxSeverity:    "Low",
			Threads:        0,
			Args:           []string{"https://www.michelin.com/"},
			ExpectedConfig: nil,
			ExpectedErr:    &internal.ErrFailedOperationOnField{"threads", "<=0", 0},
		},
		"negative-timeout": {
			SeverityFilter: "Low",
			MaxSeverity:    "Low",
			Threads:        1,
			Timeout:        -1,
			Args:           []string{"https://www.michelin.com/"},
			ExpectedConfig: nil,
			ExpectedErr:    &internal.ErrFailedOperationOnField{"timeout", "<0", -1},
		},
		"valid-config": {
			Insecure:       false,
			Export:         []string{"stdout", "csv", "json"},
			PluginFilters:  []string{},
			ExportFilename: "results",
			MaxSeverity:    "High",
			SeverityFilter: "Informational",
			URLFile:        nil,
			Threads:        1,
			Timeout:        0,
			Args:           []string{"https://www.michelin.com/"},
			ExpectedConfig: &internal.Config{
				HTTP: internal.HTTPConfig{
					Insecure: false,
					Timeout:  0,
				},
				MaxSeverity:    internal.High,
				SeverityFilter: internal.Informational,
				ExportFormats:  []string{"stdout", "csv", "json"},
				PluginFilter:   []string{},
				Urls:           []string{"https://www.michelin.com/"},
				ExportFilename: "results",
				Goroutines:     1,
			},
			ExpectedErr: nil,
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			config, err := internal.BuildConfig(tt.Insecure, tt.Export, tt.PluginFilters, tt.ExportFilename, tt.MaxSeverity, tt.SeverityFilter, tt.URLFile, tt.Threads, tt.Timeout, tt.Args)

			if !reflect.DeepEqual(config, tt.ExpectedConfig) {
				t.Errorf("Failed to get expected Config: got \"%v\" intead of \"%v\".", config, tt.ExpectedConfig)
			}

			checkErr(err, tt.ExpectedErr, t)
		})
	}
}
