package clientocol

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/inysc/routtp"
)

var reqPl = sync.Pool{
	New: func() any {
		return &Req{}
	},
}

type Client struct {
	cc         *http.Client
	headers    []routtp.Pair[string, *string]
	reqEncode  func(any) (io.Reader, string, error) // body content-type error
	respDecode func(*http.Response, any) error
	errDecode  func(*http.Response) error
}

func NewClient() *Client {
	return &Client{
		cc:         &http.Client{},
		headers:    make([]routtp.Pair[string, *string], 0, 10),
		reqEncode:  defaultReqEncode,
		respDecode: defaultRespDecode,
		errDecode:  defaultErrDecode,
	}
}

func (client *Client) WithRedirct(f func(*http.Request, []*http.Request) error) *Client {
	client.cc.CheckRedirect = f
	return client
}

func (client *Client) WithTransport(trans http.RoundTripper) *Client {
	client.cc.Transport = trans
	return client
}

func (client *Client) WithTimeout(d time.Duration) *Client {
	client.cc.Timeout = d
	return client
}

func (client *Client) WithReqEncode(f func(any) (io.Reader, string, error)) {
	client.reqEncode = f
}

func (client *Client) WithRespDecode(f func(*http.Response, any) error) {
	client.respDecode = f
}

func (client *Client) WithErrDecode(f func(*http.Response) error) {
	client.errDecode = f
}

func (client *Client) WithHeader(key string, value *string) *Client {
	client.headers = append(
		client.headers,
		routtp.Pair[string, *string]{Key: key, Val: value},
	)
	return client
}

func (client *Client) Invoke(meth string, url string, bReq any, sh ...func(*http.Header)) *Req {
	req := reqPl.Get().(*Req)
	req.cc = client

	// 序列化 请求体
	body, ctype, err := client.reqEncode(bReq)
	if err != nil {
		req.err = err
		return nil
	}

	req.req, err = http.NewRequest(meth, url, body)
	if err != nil {
		return nil
	}
	req.req.Header.Set("Content-Type", ctype)
	for _, v := range client.headers {
		if v.Val != nil {
			req.req.Header.Set(v.Key, *v.Val)
		}
	}

	for _, v := range sh {
		v(&req.req.Header)
	}

	return req
}

type Req struct {
	cc  *Client
	err error
	req *http.Request
}

func (req *Req) Do(bResp any) error {
	defer reqPl.Put(req)
	if req.err != nil {
		return req.err
	}

	resp, err := req.cc.cc.Do(req.req)
	if err != nil {
		return err
	}

	if req.cc.errDecode != nil && resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return req.cc.errDecode(resp)
	}

	return req.cc.respDecode(resp, bResp)
}

func defaultReqEncode(req any) (body io.Reader, ctype string, err error) {
	ctype = "text/plain"
	if req == nil {
		return
	}
	var bs []byte
	bs, err = json.Marshal(req)
	if err != nil {
		return
	}
	ctype = "application/json"
	body = bytes.NewReader(bs)

	return
}
func defaultRespDecode(resp *http.Response, bResp any) error {
	defer resp.Body.Close()
	if bResp == nil {
		return nil
	}
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, bResp)
}

func defaultErrDecode(resp *http.Response) error {
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return errors.New(routtp.BytesToString(bs))
}
