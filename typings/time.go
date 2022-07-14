package typings

import (
	"time"

	"github.com/inysc/ego/constant"
)

// 秒级 时间维护
type Time struct {
	sec  int64 // 毫秒
	time time.Time
}

func NewTime(sec int64) *Time {
	return &Time{
		sec:  sec,
		time: time.Unix(sec, 0),
	}
}

func (tm *Time) MarshalJSON() ([]byte, error) {
	return tm.time.AppendFormat([]byte{}, constant.TimeFmt2), nil
}

func (tm *Time) UnmarshalJSON(data []byte) error {
	return nil
}

func (tm *Time) String() string {
	return tm.Format("")
}

func (tm *Time) Format(format string) string {
	return tm.time.Format(format)
}
