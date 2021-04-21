package internal

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// SafeResults stores a Result slice. It should
// support concurrency.
type SafeResults struct {
	mux sync.Mutex
	Res []*Result
}

// Add adds a Result to the existing ones. Does not
// check for duplications.
func (s *SafeResults) Add(res *Result) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.Res = append(s.Res, res)
}

// Scanner defines the method a scanner must implement.
type Scanner interface {
	Run(urls []string, doneChan <-chan struct{}) ([]*Result, error)
}

// Scan is the entrypoint of the scan process.
func Scan(scanner Scanner, urls []string, doneChan <-chan struct{}) ([]*Result, time.Duration, error) {
	begin := time.Now()
	res, err := scanner.Run(urls, doneChan)
	if err != nil {
		return nil, time.Duration(0), err
	}

	return res, time.Since(begin), nil
}

// Scanner wraps the Signatures and the fetchers.
//
// XXX Two fetchers are needed because we can't use the same http client to follow redirects
type CoreScanner struct {
	Signatures        *Signatures
	Fetcher           Fetcher
	NoRedirectFetcher Fetcher
	Goroutines        int64
	SafeResults       *SafeResults
}

var _ = (Scanner)(&CoreScanner{})

func NewCoreScanner(config *Config, signatures *Signatures) (*CoreScanner, error) {
	// Validate parameters
	if config == nil {
		return nil, &ErrNilParameter{"config"}
	}
	if signatures == nil {
		return nil, &ErrNilParameter{"signatures"}
	}

	return &CoreScanner{
		Signatures:        signatures,
		Fetcher:           NewNetFetcher(config.HTTP.Insecure, config.HTTP.Timeout),
		NoRedirectFetcher: NewNoRedirectNetFetcher(config.HTTP.Insecure, config.HTTP.Timeout),
		SafeResults: &SafeResults{
			Res: []*Result{},
			mux: sync.Mutex{},
		},
		Goroutines: config.Goroutines,
	}, nil
}

type workerJob struct {
	url      string
	endpoint string
	plugin   Plugin
}

// RunScan scans the urls until job is completed or
// a done signal is sent throuh the chan.
func (scanner *CoreScanner) Run(urls []string, doneChan <-chan struct{}) ([]*Result, error) {
	if scanner == nil {
		return nil, &ErrNilParameter{"scanner"}
	}

	wgJobs := new(sync.WaitGroup)
	jobs := make(chan workerJob)

	// Start each job in a goroutine
	for i := int64(0); i < scanner.Goroutines; i++ {
		wgJobs.Add(1)
		go func() {
			defer wgJobs.Done()
			for {
				select {
				case <-doneChan:
					// The scan is done (force stopped)
					return
				case job, ok := <-jobs:
					// A job is here

					if !ok {
						// The job was "you do not have anymore"
						return
					}

					// Fetch the HTTP response from url
					resp, err := scanner.Fetch(job.url, job.plugin.FollowRedirects)
					if err != nil {
						logrus.Error(err)
						break
					}

					swg := new(sync.WaitGroup)
					for _, check := range job.plugin.Checks {
						swg.Add(1)
						go func(check Check) {
							defer swg.Done()
							select {
							case <-doneChan:
								return
							default:
								match, err := check.Match(resp)
								if err != nil {
									// TODO do something.
								}
								if match {
									scanner.SafeResults.Add(&Result{
										URL:         job.url,
										Name:        check.Name,
										Endpoint:    job.endpoint,
										Severity:    check.Severity,
										Remediation: check.Remediation,
									})
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
		for _, plugin := range scanner.Signatures.Plugins {
			for _, e := range plugin.Endpoints {
				endpoint := e
				if plugin.QueryString != "" {
					endpoint = fmt.Sprintf("%s?%s", endpoint, plugin.QueryString)
				}
				fullURL := fmt.Sprintf("%s%s", url, endpoint)
				logrus.Info("Testing url : ", fullURL)

				w := workerJob{url: fullURL, endpoint: endpoint, plugin: plugin}
				select {
				case <-doneChan:
					// XXX this break statement does not do anything
					break
				case jobs <- w:
				}
			}
		}
	}

	close(jobs)
	wgJobs.Wait()

	return scanner.SafeResults.Res, nil
}

// Fetch fetches content from an URL from its fetchers
// with or without redirection.
func (s CoreScanner) Fetch(url string, followRedirects bool) (*HTTPResponse, error) {
	if followRedirects {
		return Fetch(s.Fetcher, url)
	}
	return Fetch(s.NoRedirectFetcher, url)
}
