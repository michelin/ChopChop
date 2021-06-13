package internal_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/michelin/gochopchop/internal"
)

// FakeScanner mocks an internal.Scanner.
type FakeScanner struct{}

func (f *FakeScanner) Run(urls []string, doneChan <-chan struct{}) ([]internal.Result, error) {
	var res []internal.Result
	for _, url := range urls {
		switch url {
		case "https://www.michelin.com/":
			res = append(res, internal.Result{
				URL:         url,
				Endpoint:    "/",
				Name:        "michelin",
				Severity:    "Low",
				Remediation: "Work on open-sources projects.",
			})
		default:
			return nil, errFake
		}
	}
	return res, nil
}

var _ = (internal.Scanner)(&FakeScanner{})

func TestScan(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		Scanner             internal.Scanner
		URLs                []string
		DoneChan            <-chan struct{}
		ExpectedResultSlice internal.ResultSlice
		ExpectedErr         error
	}{
		"success": {
			Scanner:  &FakeScanner{},
			URLs:     []string{"https://www.michelin.com/"},
			DoneChan: nil,
			ExpectedResultSlice: internal.ResultSlice{{
				URL:         "https://www.michelin.com/",
				Endpoint:    "/",
				Name:        "michelin",
				Severity:    "Low",
				Remediation: "Work on open-sources projects.",
			}},
			ExpectedErr: nil,
		},
		"failure": {
			Scanner:             &FakeScanner{},
			URLs:                []string{""},
			DoneChan:            nil,
			ExpectedResultSlice: nil,
			ExpectedErr:         errFake,
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			resp, _, err := internal.Scan(tt.Scanner, tt.URLs, tt.DoneChan)

			if !cmp.Equal(resp, tt.ExpectedResultSlice) {
				t.Errorf("Failed to get expected ResultSlice: got \"%v\" instead of \"%v\".", resp, tt.ExpectedResultSlice)
			}
			checkErr(err, tt.ExpectedErr, t)
		})
	}
}

func TestNewCoreScanner(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		Config              *internal.Config
		Signatures          *internal.Signatures
		ExpectedCoreScanner *internal.CoreScanner
		ExpectedErr         error
	}{
		"nil-config": {
			Config:              nil,
			Signatures:          nil,
			ExpectedCoreScanner: nil,
			ExpectedErr:         &internal.ErrNilParameter{"config"},
		},
		"nil-signatures": {
			Config:              &internal.Config{},
			Signatures:          nil,
			ExpectedCoreScanner: nil,
			ExpectedErr:         &internal.ErrNilParameter{"signatures"},
		},
		"core-scanner": {
			Config:              &internal.Config{},
			Signatures:          &internal.Signatures{},
			ExpectedCoreScanner: &internal.CoreScanner{},
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			scan, err := internal.NewCoreScanner(tt.Config, tt.Signatures)

			if (scan == nil) != (tt.ExpectedCoreScanner == nil) {
				t.Errorf("Failed to get a non-nil CoreScanner: got \"%v\" instead of \"%v\".", scan, tt.ExpectedCoreScanner)
			}
			checkErr(err, tt.ExpectedErr, t)
		})
	}
}

func TestCodeScannerFetch(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		CoreScanner     internal.CoreScanner
		URL             string
		FollowRedirects bool
		ExpectedResp    *internal.HTTPResponse
		ExpectedErr     error
	}{}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			resp, err := tt.CoreScanner.Fetch(tt.URL, tt.FollowRedirects)
			if !reflect.DeepEqual(resp, tt.ExpectedResp) {
				t.Errorf("Failed to get expected HTTPResponse: got \"%v\" instead of \"%v\".", resp, tt.ExpectedResp)
			}

			checkErr(err, tt.ExpectedErr, t)
		})
	}
}
