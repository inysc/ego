package clientocol

import (
	"net/http"
	"sync"
)

var reqPl = sync.Pool{
	New: func() any {
		return &Req{}
	},
}

type Client struct {
	cc   http.Client
	resp func(*http.Response, any) error
	err  func(*http.Response, any) error
}

func (client *Client) Invoke(meth string) *Req {
	req := reqPl.Get().(*Req)
	req.cc = &client.cc

	return req
}

type Req struct {
	cc   *http.Client
	req  *http.Request
	err  func(*http.Response) error
	resp func(*http.Response, any) error
}

func (req *Req) Do(bResp any) error {
	defer reqPl.Put(req)
	resp, err := req.cc.Do(req.req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return req.err(resp)
	}

	return req.resp(resp, bResp)
}
