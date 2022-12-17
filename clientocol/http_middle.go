package clientocol

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/inysc/facade"
	"github.com/inysc/routtp"
)

const (
	HeaderQiuUser  = "X-Qiu-User"
	HeaderQiuPrior = "X-Qiu-Prior"
	HeaderQiuCity  = "X-Qiu-City"
	HeaderQiuIp    = "X-Qiu-Ip"
	HeaderQiuHost  = "X-Qiu-Host"
)

// Logger 接收 routtp 框架默认的日志
func Logger() routtp.Handler {
	return func(ctx *routtp.Context) {
		start := time.Now()
		ctx.Next()

		facade.Infof(
			"%s method[%s] query[%s] ip[%s] userAgent[%s] cost[%s]",
			ctx.Request.URL.Path, ctx.Request.Method, ctx.Request.URL.RawQuery,
			ctx.RealIP(), ctx.Request.UserAgent(), time.Since(start),
		)
	}
}

// Recovery recover 掉项目可能出现的 panic
func Recovery() routtp.Handler {
	return func(ctx *routtp.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(ctx.Request, false)
				if brokenPipe {
					path := ctx.Request.URL.Path
					facade.Errorf("path[%s] error[%s] request[%s]", path, err, httpRequest)

					ctx.Exception(&exception{
						code:   400,
						status: 400,
						msg:    err.(error).Error(),
						data:   "",
					})
					ctx.Abort()
					return
				}

				facade.Errorf("[Recovery from panic], err[%s], request[%s], stack\n%s",
					err, httpRequest, debug.Stack())
				ctx.Response.WriteHeader(http.StatusInternalServerError)
			}
		}()
		ctx.Next()
	}
}

type exception struct {
	code   int
	status int
	msg    string
	data   string
}

func (ep *exception) Code() int {
	return ep.code
}

func (ep *exception) Msg() string {
	return ep.msg
}

func (ep *exception) Data() any {
	return ep.data
}

func (ep *exception) Status() int {
	return ep.status
}
