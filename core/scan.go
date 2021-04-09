package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/michelin/gochopchop/internal"

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
	Threads  int
}

// NewScanner returns a pointer to a initialized Scanner
func NewScanner(fetcher IFetcher, noRedirectFetcher IFetcher, signatures *Signatures, threads int) *Scanner {
	safeData := &SafeData{out: make([]Output, 0)}
	return &Scanner{
		Signatures:        signatures,
		Fetcher:           fetcher,
		NoRedirectFetcher: noRedirectFetcher,
		safeData:          safeData,
		Threads:           threads,
	}
}

type workerJob struct {
	url      string
	endpoint string
	plugin   *Plugin
}

func (s Scanner) Scan(ctx context.Context, urls []string) ([]Output, error) {
	wg := new(sync.WaitGroup)
	jobs := make(chan workerJob)

	for i := 0; i < s.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobs:
					if !ok { // no more jobs
						return
					}
					resp, err := s.fetch(job.url, job.plugin.FollowRedirects)
					if err != nil {
						log.Error(err)
						break
					}
					swg := new(sync.WaitGroup)
					for _, check := range job.plugin.Checks {
						swg.Add(1)
						go func(check *Check) {
							defer swg.Done()
							select {
							case <-ctx.Done():
								return
							default:
								if check.Match(resp) {
									o := Output{
										URL:         job.url,
										Name:        check.Name,
										Endpoint:    job.endpoint,
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
			}
		}()
	}

	for _, url := range urls {
		for _, plugin := range s.Signatures.Plugins {
			if plugin.Endpoint != "" {
				plugin.Endpoints = []string{plugin.Endpoint}
			}
			for _, e := range plugin.Endpoints {
				endpoint := e
				if plugin.QueryString != "" {
					endpoint = fmt.Sprintf("%s?%s", endpoint, plugin.QueryString)
				}
				fullURL := fmt.Sprintf("%s%s", url, endpoint)
				log.Info("Testing url : ", fullURL)

				w := workerJob{url: fullURL, endpoint: endpoint, plugin: plugin}
				select {
				case <-ctx.Done():
					break
				case jobs <- w:
				}
			}
		}
	}

	close(jobs)
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
	return httpResponse, nil
}
