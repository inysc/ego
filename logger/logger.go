package logger

import (
	"sync"

	"github.com/inysc/qog"
)

var (
	l *qog.Logger
	o sync.Once
)

func GetLogger(srvname, filename string) *qog.Logger {
	o.Do(func() {
		if l == nil {
			l = qog.New(srvname, qog.DEBUG, &Logger{
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
