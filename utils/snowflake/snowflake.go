package snowflake

import (
	"sync"
	"time"
)

type snowFlake struct {
	y   int
	ts  uint64
	seq uint64
	mut sync.Mutex
}

const (
	epoch  = uint64(1650025543467663)    // 设置起始时间
	bitSeq = uint(6)                     // 序列所占的位数
	bitTS  = uint(64 - 6)                // 时间戳占用位数
	maxSeq = uint64(-1 ^ (-1 << bitSeq)) // 支持的最大序列id数量
	maxTS  = uint64(-1 ^ (-1 << bitTS))  // 时间戳最大值
)

var sf *snowFlake = nil

func init() {
	if sf == nil {
		sf = &snowFlake{}
	}
}

func GetVal() (val uint64) {
	sf.mut.Lock()

	now := time.Now()
	if now.Year() != sf.y || sf.seq >= maxSeq {
		sf.y = now.Year()
		sf.ts = uint64(now.UnixMicro())
		sf.seq = 0
	}
	val = (sf.ts-epoch)<<bitSeq + sf.seq

	sf.seq++
	sf.mut.Unlock()
	return
}
