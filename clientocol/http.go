package clientocol

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/inysc/facade"
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
		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(bs, respBody)
		return err
	}

	return nil
}

func Start(srv *http.Server) {
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			facade.Errorf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	facade.Infof("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		facade.Errorf("server shutdown[%s]", err)
	}

	facade.Infof("Server exiting")
}
