package internal

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Scanner defines the method a scanner must implement.
type Scanner interface {
	Run(urls []string, doneChan <-chan struct{}) ([]Result, error)
}

// Scan is the entrypoint of the scan process.
func Scan(scanner Scanner, urls []string, doneChan <-chan struct{}) (ResultSlice, time.Duration, error) {
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
	Goroutines        int
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
			Res: []Result{},
			mux: sync.Mutex{},
		},
		Goroutines: config.Goroutines,
	}, nil
}

type workerJob struct {
	url      string
	endpoint string
	plugin   *Plugin
}

// RunScan scans the urls until job is completed or
// a done signal is sent throuh the chan.
func (scanner *CoreScanner) Run(urls []string, doneChan <-chan struct{}) ([]Result, error) {
	// Split the load in channels to work on concurrently
	workJobs := splitWork(scanner.Goroutines, urls, scanner.Signatures.Plugins)
	workJobsChan := channelize(workJobs)

	var wgJobs sync.WaitGroup
	for _, wj := range workJobsChan {
		wgJobs.Add(1)
		go func(wj chan workerJob) {
			defer wgJobs.Done()
			for {
				select {
				case <-doneChan:
					return

				default:
					return

				case job := <-wj:
					// Fetch the HTTP response from url
					resp, err := scanner.Fetch(job.url+job.endpoint, job.plugin.FollowRedirects)
					if err != nil {
						logrus.Error(err)
						break
					}

					// Procede to checks concurrently
					var swg sync.WaitGroup
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
									return
								}
								if match {
									scanner.SafeResults.Append(Result{
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
		}(wj)
	}

	wgJobs.Wait()
	for _, wj := range workJobsChan {
		close(wj)
	}

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

// This is the implementation of the following algorithm, which aims to
// split the work into pieces, in near-equal parts, with a bit of optimization
// to avoid working on empty stuff.
//
// Let's have a set to split in n parts. How to tend to a prefectly-splitted
// load from this ? Dividing it seems to be a good solution.
//
// In our case, we are having endpoints and threads (it will be implemented on
// goroutines instead of threads, so only keep the idea behind).
// \#endpoints = \#urls \times \sum_{p=0}^{\#plugins} \#endpoints_p
//
// To split the entrypoints to subsets (integer-long), we divide it by our number
// of subsets, e.g. number of threads, which gives us the length of a "standard"
// subset. The remaining goes into a last subset.
// L_s stands for "standard length" ; L_r stands for "remaining length".
// L_s = \left \lceil \frac{\#endpoints}{\#threads} \right \rceil
// L_r = \#endpoints - L_s \times (\#threads - 1)
// Notice L_s \geq L_r.
//
// We now need to define how many subset of length L_s and L_r will be built,
// and by construction we know the following.
// \#L_s = \#threads - 1
// \#L_r = 1
//
// Examples:
// - Case 16/4:
//   \#endpoints = 16 ; \#threads = 4
//   L_s = ceil(16/4) = 4
//   L_r = 16 - 4*(4-1) = 4
//   We can represent it as follows.
//   ■ ■ ■ ■
//   ■ ■ ■ ■
//   ■ ■ ■ ■
//   ■ ■ ■ ■
//   0 1 2 3
// - Case 21/4:
//   \#endpoints = 21 ; \#threads = 4
//   L_s = ceil(21/4) = 6
//   L_r = 21 - 6*(4-1) = 3
//   We can represent it as follows.
//   ■ ■ ■
//   ■ ■ ■
//   ■ ■ ■
//   ■ ■ ■ ■
//   ■ ■ ■ ■
//   ■ ■ ■ ■
//   0 1 2 3
//
// We can see that the last one is taking a nap while the other are still
// working.
// nap = L_s - L_r
// So, we want to move "toFill" endpoints to the last subset.
// toFill = nap - 1 = L_s - L_r - 1
//
// Let's change our concept and talk of "filled" and "reduced" subsets.
// To avoid having the last taking a nap while the other are working, we "move" last
// endpoint of toFill "standard" subsets to this last one, so we have the following.
// \#reduced = toFill + 1
// \#filled = \#threads - \#reduced = \#threads - (toFill + 1) = \#threads - toFill - 1
//
// Then, by construction, we define L_f (filled length) and L_r (reduced length) as follows.
// L_f = \left \lceil \frac{\#endpoints}{\#threads} \right \rceil
// L_r = \left \lfloor \frac{\#endpoints}{\#threads} \right \rfloor
//
// Now, to avoid relaying on the previous step to achieve this one, we develop and reduce.
// \#reduced = toFill + 1
//           = L_s - L_r - 1                        (L_r from the previous algorithm)
//           = \left \lceil \frac{\#endpoints}{\#threads} \right \rceil - (\#endpoints - \left \lceil \frac{\#endpoints}{\#threads} \right \rceil \times (\#threads - 1))
//           = \left \lceil \frac{\#endpoints}{\#threads} \right \rceil - \#endpoints + \left \lceil \frac{\#endpoints}{\#threads} \right \rceil \times (\#threads - 1)
//           = \left \lceil \frac{\#endpoints}{\#threads} \right \rceil \times \#threads - \#endpoints
//
// Examples:
// - Case 16/4:
//   \#endpoints = 16 ; \#threads = 4
//   L_f = ceil(16/4) = 4
//   L_r = floor(16/4) = 4
//   #reduced = 4 * 4 - 16 = 0
//   #filled = 4 - 0 = 4
//   We can represent it as follows.
//   ■ ■ ■ ■
//   ■ ■ ■ ■
//   ■ ■ ■ ■
//   ■ ■ ■ ■
//   0 1 2 3
// - Case 21/4:
//   \#endpoints = 21 ; \#threads = 4
//   L_f = ceil(21/4) = 6
//   L_r = floor(21/4) = 5
//   #reduced = 6 * 4 - 21 = 3
//   #filled = 4 - 3 = 1
//   We can represent it as follows.
//   ■
//   ■ ■ ■ ■
//   ■ ■ ■ ■
//   ■ ■ ■ ■
//   ■ ■ ■ ■
//   ■ ■ ■ ■
//   0 1 2 3
//
// For the first version, in case L_s == \#threads - 1 (worst case) this last
// optimization avoids the last worker to work only on one job (see it with
// the case 157/14 where 13 works on 12 jobs and the last on only one).
func splitWork(pieces int, urls []string, plugins []Plugin) [][]workerJob {
	// Compute #endpoints
	var n_endpoints int
	for _, plugin := range plugins {
		n_endpoints += len(plugin.Endpoints)
	}
	n_endpoints *= len(urls)

	// Compute L_filled and L_reduced
	l_filled := n_endpoints / pieces
	l_reduced := l_filled
	if l_filled*pieces != n_endpoints {
		// Ceil l_filled
		l_filled++
	}

	// Compute #reduced and #filled
	n_reduced := l_filled*pieces - n_endpoints
	n_filled := pieces - n_reduced

	// Avoid having empty work groups
	if l_filled == 0 {
		return [][]workerJob{}
	}
	if l_reduced == 0 {
		n_reduced = 0
	}

	// Build work groups
	n_total := n_filled + n_reduced
	wg := make([][]workerJob, n_total)
	for i := 0; i < n_filled; i++ {
		wg[i] = make([]workerJob, l_filled)
	}
	for i := 0; i < n_reduced; i++ {
		wg[n_filled+i] = make([]workerJob, l_reduced)
	}

	// Load balance
	indexWG := 0          // Index of the work group
	indexInCurrWG := 0    // Index in the current work group
	lenCurrWg := l_filled // Length of the current work group
	currWG := &wg[0]      // Place iterator at first work group
	lp := len(plugins)
	for _, url := range urls {
		for i := 0; i < lp; i++ {
			for _, endp := range plugins[i].Endpoints {
				// Set the work group endpoint
				(*currWG)[indexInCurrWG] = workerJob{
					url:      url,
					endpoint: endp,
					plugin:   &plugins[i],
				}
				indexInCurrWG++

				// Move the indexes if needed.
				if indexInCurrWG == lenCurrWg {
					indexInCurrWG = 0
					indexWG++

					if indexWG == n_filled {
						lenCurrWg = l_reduced
					}

					if indexWG != n_total {
						// Move to the next only if not the last in slice
						currWG = &wg[indexWG]
					}
				}
			}
		}
	}

	return wg
}

func channelize(workerJobs [][]workerJob) []chan workerJob {
	chans := make([]chan workerJob, len(workerJobs))
	for i, wj := range workerJobs {
		chans[i] = make(chan workerJob, len(wj))
		for _, w := range wj {
			chans[i] <- w
		}
	}

	return chans
}
