package pkg

import (
    "crypto/tls"
    "log"
    "net/http"
    "time"
)

//HTTPGet return http response of http get request
func HTTPGet(insecure bool, url string) (*http.Response, error) {
    tr := &http.Transport{}
    if insecure {
        tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
    }
    var netClient = &http.Client{
        Transport: tr,
        Timeout:   time.Second * 3,
    }
    resp, err := netClient.Get(url)
    if err != nil {
        log.Println(err)
        log.Println("If error unsupported protocol scheme encountered, try adding flag --prefix with http://, or add prefix directly in url list")
    }

    return resp, err
}
