package internal

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
)

// IHTTPClient is the interface defining a struct that
// get a *http.Response from an URL.
type IHTTPClient interface {
	Get(url string) (*http.Response, error)
}

// HTTPClient wraps the HTTP Client parameters to use.
type HTTPClient struct {
	Transport http.RoundTripper
	Timeout   time.Duration
}

// HTTPResponse wraps the status code, body and header
// of a http.Response (it is a simplified version of
// this last).
type HTTPResponse struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

// Fetcher wraps an IHTTPClient and defines a method to
// fetch an URL.
type Fetcher struct {
	Netclient IHTTPClient
}

// NewFetcher builds and returns a new fetcher.
func NewFetcher(insecure bool, timeout int64) *Fetcher {
	var tlsCfg *tls.Config
	if insecure {
		tlsCfg = &tls.Config{InsecureSkipVerify: true}
	}

	return &Fetcher{
		Netclient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsCfg,
			},
			Timeout: time.Second * time.Duration(timeout),
		},
	}
}

// NewNoRedirectFetcher builds and returns a new
// fetcher that does not redirect.
func NewNoRedirectFetcher(insecure bool, timeout int64) *Fetcher {
	var tlsCfg *tls.Config
	if insecure {
		tlsCfg = &tls.Config{InsecureSkipVerify: true}
	}

	return &Fetcher{
		Netclient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsCfg,
			},
			Timeout: time.Second * time.Duration(timeout),
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

// Fetch fetches a *http.Response from an URL and returns
// a *HTTPResponse from it.
func (f Fetcher) Fetch(url string) (*HTTPResponse, error) {
	resp, err := f.Netclient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &HTTPResponse{
		Body:       body,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
	}, nil
}
