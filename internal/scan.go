package internal

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// SafeData stores a Result slice. It should
// support concurrency.
type SafeData struct {
	mux sync.Mutex
	res []*Result
}

// Add adds a Result to the existing ones. Does not
// check for duplications.
func (s *SafeData) Add(res *Result) {
	s.mux.Lock()
	s.res = append(s.res, res)
	s.mux.Unlock()
}

// IFetcher defines the method to fetch a HTTP response
// from an URL.
type IFetcher interface {
	Fetch(url string) (*HTTPResponse, error)
}

// IScanner defines the method to fetch Results from a slice
// of URLs.
type IScanner interface {
	Scan(urls []string) ([]Result, error)
}

// Scanner wraps the Signatures and the fetchers.
//
// TODO refactor this shit...
// XXX Two fetchers are needed because we can't use the same http client to follow redirects
type Scanner struct {
	Signatures        *Signatures
	Fetcher           IFetcher
	NoRedirectFetcher IFetcher
	safeData          *SafeData
	Goroutines        int64
}

type workerJob struct {
	url      string
	endpoint string
	plugin   *Plugin
}

// RunScan scans the urls until job is completed or
// a done signal is sent throuh the chan.
func RunScan(scanner *Scanner, urls []string, doneChan <-chan struct{}) ([]*Result, error) {
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
						go func(check *Check) {
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
									scanner.safeData.Add(&Result{
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

	return scanner.safeData.res, nil
}

// Fetch fetches content from an URL from its fetchers
// with or without redirection.
func (s Scanner) Fetch(url string, followRedirects bool) (*HTTPResponse, error) {
	var httpResponse *HTTPResponse
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

// Scan scans a set of URLs through the given config and
// builds the results according to the given signatures.
// It stops execution on a signal through "doneChan".
//
// It returns the results and the duration of the scan.
func Scan(config *Config, signatures *Signatures, doneChan <-chan struct{}) ([]*Result, time.Duration, error) {
	// Validate parameters
	if config == nil {
		return nil, time.Duration(0), &ErrNilParameter{"config"}
	}
	if signatures == nil {
		return nil, time.Duration(0), &ErrNilParameter{"signatures"}
	}

	// Build fetchers
	fetcher := NewFetcher(config.HTTP.Insecure, config.HTTP.Timeout)
	noRdrFetcher := NewNoRedirectFetcher(config.HTTP.Insecure, config.HTTP.Timeout)

	// Run the scan
	scanner := &Scanner{
		Signatures:        signatures,
		Fetcher:           fetcher,
		NoRedirectFetcher: noRdrFetcher,
		safeData: &SafeData{
			res: []*Result{},
			mux: sync.Mutex{},
		},
		Goroutines: config.Goroutines,
	}
	begin := time.Now()
	res, err := RunScan(scanner, config.Urls, doneChan)
	if err != nil {
		return nil, time.Duration(0), err
	}

	return res, time.Since(begin), nil
}
