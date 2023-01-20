package httpget

import (
	"crypto/tls"
	"gochopchop/internal"
	"gochopchop/core"
	"io/ioutil"
	"net/http"
	"time"
)

type IHTTPClient interface {
	Get(url string) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
}

type HTTPClient struct {
	Transport http.RoundTripper
	Timeout   time.Duration
}

type Fetcher struct {
	Netclient IHTTPClient
	Config *core.Config
}

func NewFetcher(config *core.Config) *Fetcher {
	tr := &http.Transport{}
	if config.HTTP.Insecure {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	var netClient = &http.Client{
		Transport: tr,
		Timeout:   time.Second * time.Duration(config.HTTP.Timeout),
	}
	return &Fetcher{
		Netclient: netClient,
		Config: config,
	}
}

func NewNoRedirectFetcher(config *core.Config) *Fetcher {
	tr := &http.Transport{}
	if config.HTTP.Insecure {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	var netClient = &http.Client{
		Transport: tr,
		Timeout:   time.Second * time.Duration(config.HTTP.Timeout),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return &Fetcher{
		Netclient: netClient,
		Config: config,
	}
}

func (s Fetcher) Fetch(url string) (*internal.HTTPResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", s.Config.HTTP.UserAgent)

	resp, err := s.Netclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)

	var r = &internal.HTTPResponse{
		Body:       bodyString,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
	}

	return r, err
}
