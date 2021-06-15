package internal_test

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/michelin/gochopchop/internal"
)

func TestNewFetcher(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		Insecure        bool
		Timeout         int
		ExpectedFetcher internal.NetFetcher
	}{
		"secure": {
			Insecure: false,
			Timeout:  0,
			ExpectedFetcher: internal.NetFetcher{
				Client: http.Client{
					Transport: &http.Transport{
						TLSClientConfig: nil,
					},
					Timeout: time.Duration(0),
				},
			},
		},
		"insecure": {
			Insecure: true,
			Timeout:  0,
			ExpectedFetcher: internal.NetFetcher{
				Client: http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					},
					Timeout: time.Duration(0),
				},
			},
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			f := internal.NewNetFetcher(tt.Insecure, tt.Timeout)

			if tt.Insecure != (f.Client.Transport.(*http.Transport).TLSClientConfig != nil) {
				t.Errorf("Failed to get a TLS config appropriate with the insecure parameter (%t): \"%v\".", tt.Insecure, f)
			}
			if time.Second*time.Duration(tt.Timeout) != f.Client.Timeout {
				t.Errorf("Failed to get expected timeout: got \"%v\" instead of \"%v\".", f.Client.Timeout, time.Second*time.Duration(tt.Timeout))
			}
		})
	}
}

func TestNewNoRedirectNetFetcher(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		Insecure        bool
		Timeout         int
		ExpectedFetcher internal.NetFetcher
	}{
		"secure": {
			Insecure: false,
			Timeout:  0,
			ExpectedFetcher: internal.NetFetcher{
				Client: http.Client{
					Transport: &http.Transport{
						TLSClientConfig: nil,
					},
					Timeout: time.Duration(0),
				},
			},
		},
		"insecure": {
			Insecure: true,
			Timeout:  0,
			ExpectedFetcher: internal.NetFetcher{
				Client: http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					},
					Timeout: time.Duration(0),
				},
			},
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			f := internal.NewNoRedirectNetFetcher(tt.Insecure, tt.Timeout)

			if tt.Insecure != (f.Client.Transport.(*http.Transport).TLSClientConfig != nil) {
				t.Errorf("Failed to get a TLS config appropriate with the insecure parameter (%t): \"%v\".", tt.Insecure, f)
			}
			if time.Second*time.Duration(tt.Timeout) != f.Client.Timeout {
				t.Errorf("Failed to get expected timeout: got \"%v\" instead of \"%v\".", f.Client.Timeout, time.Second*time.Duration(tt.Timeout))
			}
			if f.Client.CheckRedirect(nil, nil) != http.ErrUseLastResponse {
				t.Error("Failed to set the no-redirection method properly.")
			}
		})
	}
}

var michelinHTTPResponse = &internal.HTTPResponse{
	Body:       []byte("gochopchop\n"),
	StatusCode: 200,
	Header: http.Header{
		"Fake-Header": []string{"fake-content"},
	},
}

var errFake = errors.New("fake error")

// FailingReadCloser implements the io.Reader and io.Closer
// interfaces. It will return an error on it's Read method
// call.
type FailingReadCloser struct{}

func (f *FailingReadCloser) Read([]byte) (int, error) {
	return 0, errFake
}

func (f *FailingReadCloser) Close() error {
	return nil
}

var _ = (io.ReadCloser)(&FailingReadCloser{})

// FakeFetcher mocks a Fetcher for the following test.
type FakeFetcher struct{}

func (f FakeFetcher) Get(url string) (*http.Response, error) {
	switch url {
	case "https://www.michelin.com/":
		return &http.Response{
			Body:       NewFakeReadCloser("gochopchop\n"),
			StatusCode: 200,
			Header: http.Header{
				"Fake-Header": []string{"fake-content"},
			},
		}, nil
	case "gochopchop":
		return &http.Response{
			Body: &FailingReadCloser{},
		}, nil
	default:
		return nil, errFake
	}
}

func TestFetch(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		Fetcher      internal.Fetcher
		URL          string
		ExpectedResp *internal.HTTPResponse
		ExpectedErr  error
	}{
		"valid-url": {
			Fetcher:      FakeFetcher{},
			URL:          "https://www.michelin.com/",
			ExpectedResp: michelinHTTPResponse,
			ExpectedErr:  nil,
		},
		"fail-read": {
			Fetcher:      FakeFetcher{},
			URL:          "gochopchop",
			ExpectedResp: nil,
			ExpectedErr:  errFake,
		},
		"invalid-url": {
			Fetcher:      FakeFetcher{},
			URL:          "",
			ExpectedResp: nil,
			ExpectedErr:  errFake,
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			resp, err := internal.Fetch(tt.Fetcher, tt.URL)

			if !reflect.DeepEqual(resp, tt.ExpectedResp) {
				t.Errorf("Failed to get expected HTTPResponse: got \"%v\" instead of \"%v\".", resp, tt.ExpectedResp)
			}
			checkErr(err, tt.ExpectedErr, t)
		})
	}
}
