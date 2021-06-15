package internal

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

// HTTPResponse wraps the status code, body and header
// of a http.Response (it is a simplified version of
// this last).
type HTTPResponse struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

// Fetcher defines the method a http fetcher should define.
type Fetcher interface {
	// Get takes an URL and returns a *http.Response from it.
	Get(url string) (*http.Response, error)
}

// Fetcher wraps an IHTTPClient and defines a method to
// fetch an URL.
type NetFetcher struct {
	Client http.Client
}

// NewNetFetcher builds and returns a new fetcher.
func NewNetFetcher(insecure bool, timeout int) NetFetcher {
	var tlsCfg *tls.Config
	if insecure {
		tlsCfg = &tls.Config{InsecureSkipVerify: true}
	}

	return NetFetcher{
		Client: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsCfg,
			},
			Timeout: time.Second * time.Duration(timeout),
		},
	}
}

// NewNoRedirectFetcher builds and returns a new
// fetcher that does not follow redirection.
func NewNoRedirectNetFetcher(insecure bool, timeout int) NetFetcher {
	var tlsCfg *tls.Config
	if insecure {
		tlsCfg = &tls.Config{InsecureSkipVerify: true}
	}

	return NetFetcher{
		Client: http.Client{
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

func (nf NetFetcher) Get(url string) (*http.Response, error) {
	return nf.Client.Get(url)
}

// Fetch takes a Fetcher and an URL to return a *HTTPResponse from.
func Fetch(fetcher Fetcher, url string) (*HTTPResponse, error) {
	resp, err := fetcher.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &HTTPResponse{
		Body:       body,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
	}, nil
}
