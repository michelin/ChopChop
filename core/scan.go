package core

import (
	"context"
	"fmt"
	"gochopchop/internal"
	"sync"

	log "github.com/sirupsen/logrus"
)

type SafeData struct {
	mux sync.Mutex
	out []Output
}

func (s *SafeData) Add(d Output) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.out = append(s.out, d)
}

type IFetcher interface {
	Fetch(url string) (*internal.HTTPResponse, error)
}

type IScanner interface {
	Scan(urls []string) ([]Output, error)
}

type Scanner struct {
	Signatures        *Signatures
	Fetcher           IFetcher
	NoRedirectFetcher IFetcher
	// Two fetchers are needed because we can't use the same http client to follow redirects
	safeData *SafeData
}

func NewScanner(fetcher IFetcher, noRedirectFetcher IFetcher, signatures *Signatures) *Scanner {
	safeData := new(SafeData)
	return &Scanner{
		Signatures:        signatures,
		Fetcher:           fetcher,
		NoRedirectFetcher: noRedirectFetcher,
		safeData:          safeData,
	}
}

func (s Scanner) Scan(ctx context.Context, urls []string) ([]Output, error) {
	wg := new(sync.WaitGroup)

	for _, url := range urls {
		log.Info("Testing url : ", url)
		for _, plugin := range s.Signatures.Plugins {
			for _, uri := range plugin.URIs {
				endpoint := uri
				if plugin.QueryString != "" {
					endpoint = fmt.Sprintf("%s?%s", endpoint, plugin.QueryString)
				}
				fullURL := fmt.Sprintf("%s%s", url, endpoint)
				log.Info("Testing full url : ", fullURL)

				wg.Add(1)
				go func(plugin *Plugin) {
					defer wg.Done()
					select {
					case <-ctx.Done():
						return
					default:
						resp, err := s.fetch(fullURL, plugin.FollowRedirects)
						if err != nil {
							log.Debug(err)
							return
						}
						swg := new(sync.WaitGroup)
						for _, check := range plugin.Checks {
							swg.Add(1)
							go func(check *Check) {
								defer swg.Done()
								select {
								case <-ctx.Done():
									return
								default:
									if check.Match(resp) {
										o := Output{
											URL:         url,
											Name:        check.Name,
											Endpoint:    endpoint,
											Severity:    check.Severity,
											Remediation: check.Remediation,
										}
										s.safeData.Add(o)
									}
								}
							}(check)
						}
						swg.Wait()
					}
				}(plugin)
			}
		}
	}
	wg.Wait()

	return s.safeData.out, nil
}

func (s Scanner) fetch(url string, followRedirects bool) (*internal.HTTPResponse, error) {
	var httpResponse *internal.HTTPResponse
	var err error

	if !followRedirects {
		httpResponse, err = s.NoRedirectFetcher.Fetch(url)
	} else {
		httpResponse, err = s.Fetcher.Fetch(url)
	}
	if err != nil {
		return nil, err
	}
	// weird case when both the error and the response are nil, caused by the server refusing the connection
	if httpResponse == nil {
		return nil, fmt.Errorf("Server refused the connection for : %s", url)
	}
	return httpResponse, nil
}
