package logger

import (
	"sync"

	"github.com/inysc/qog"
)

type log interface {
	Trace(string)
	Tracef(string, ...any)
	Debug(string)
	Debugf(string, ...any)
	Info(string)
	Infof(string, ...any)
	Warn(string)
	Warnf(string, ...any)
	Error(string)
	Errorf(string, ...any)
}

var (
	o sync.Once
	l log = nopLogger{}
)

func SetLogger(lg log) {
	l = lg
}

func GetLogger(srvname, filename string) log {
	o.Do(func() {
		if l == nil {
			l = qog.New(srvname, qog.DEBUG, &qog.LoggerFile{
				Filename:   filename,
				MaxSize:    30,
				MaxAge:     30,
				MaxBackups: 7,
				LocalTime:  true,
				Compress:   true,
			})
		}
	})
	return l
}

type nopLogger struct{}

func (l nopLogger) Trace(string)          {}
func (l nopLogger) Tracef(string, ...any) {}
func (l nopLogger) Debug(string)          {}
func (l nopLogger) Debugf(string, ...any) {}
func (l nopLogger) Info(string)           {}
func (l nopLogger) Infof(string, ...any)  {}
func (l nopLogger) Warn(string)           {}
func (l nopLogger) Warnf(string, ...any)  {}
func (l nopLogger) Error(string)          {}
func (l nopLogger) Errorf(string, ...any) {}
