package client

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var hcpl = sync.Pool{
	New: func() any {
		return &http.Client{
			Timeout: time.Minute,
		}
	},
}

func RawHTTPRequest(req *http.Request) (*http.Response, error) {
	client := hcpl.Get().(*http.Client)
	defer hcpl.Put(client)

	client.Transport = nil
	if req.URL.Scheme == "https" {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	return client.Do(req)
}

func HTTPRequest[T any](req *http.Request, respBody *T) error {
	resp, err := RawHTTPRequest(req)
	if err != nil {
		return err
	}
	if respBody != nil {
		defer resp.Body.Close()
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(bs, respBody)
	}

	return err
}
