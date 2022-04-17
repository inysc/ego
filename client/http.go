package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

func Start(srv http.Server, logs logger) {
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logs.Errorf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logs.Infof("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		logs.Errorf("server shutdown[%s]", err)
	}

	logs.Infof("Server exiting")
}
