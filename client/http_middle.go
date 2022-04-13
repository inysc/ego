package client

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	AuthUser  = "Utluser"
	AuthPrior = "Utlprior"
	AuthCity  = "Utlcity"
	AuthIP    = "Utlip"
)

// GinLogger 接收 gin 框架默认的日志
func GinLogger(log logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		log.Infof("%s, status[%d], method[%s], path[%s],"+
			" query[%s], ip[%s], userAgent[%s], errors[%s], cost[%s]",
			c.Writer.Status(), c.Request.Method, path, query,
			c.ClientIP(), c.Request.UserAgent(),
			c.Errors.ByType(gin.ErrorTypePrivate).String(), cost)
	}
}

// GinRecovery recover 掉项目可能出现的 panic
func GinRecovery(log logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					log.Errorf("path[%s], error[%s], request[%s]",
						c.Request.URL.Path, err, httpRequest)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				log.Errorf("[Recovery from panic], err[%s], request[%s], stack\n%s",
					err, httpRequest, debug.Stack())
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
