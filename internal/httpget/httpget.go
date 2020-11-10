package httpget

import (
	"crypto/tls"
	"gochopchop/internal"
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

type Fetcher struct {
	Netclient IHTTPClient
}

func NewFetcher(insecure bool, timeout int) *Fetcher {
	tr := &http.Transport{}
	if insecure {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	var netClient = &http.Client{
		Transport: tr,
		Timeout:   time.Second * time.Duration(timeout),
	}
	return &Fetcher{
		Netclient: netClient,
	}
}

func NewNoRedirectFetcher(insecure bool, timeout int) *Fetcher {
	tr := &http.Transport{}
	if insecure {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	var netClient = &http.Client{
		Transport: tr,
		Timeout:   time.Second * time.Duration(timeout),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return &Fetcher{
		Netclient: netClient,
	}
}

func (s Fetcher) Fetch(url string) (*internal.HTTPResponse, error) {

	resp, err := s.Netclient.Get(url)
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
