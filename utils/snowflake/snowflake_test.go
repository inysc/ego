package snowflake

import (
	"testing"
	"time"
)

func TestSF(t *testing.T) {

	t.Log(time.Now().Unix())
	t.Log(time.Now().UnixMilli())
	t.Log(time.Now().UnixMicro())
	v1 := GetVal()
	t.Log(v1)
	for range [100000000]struct{}{} {
		v2 := GetVal()
		if v1 >= v2 {
			t.Log(v1, v2)
			panic("asd")
		}
	}
}
