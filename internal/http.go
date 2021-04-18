package internal

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
)

type IHTTPClient interface {
	Get(url string) (*http.Response, error)
}

type HTTPClient struct {
	Transport http.RoundTripper
	Timeout   time.Duration
}

type HTTPResponse struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

type Fetcher struct {
	Netclient IHTTPClient
}

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
