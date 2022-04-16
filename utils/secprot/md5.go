package secprot

import (
	"crypto/md5"

	"github.com/inysc/ego/utils/bytestr"
)

func Md5Encode(val string) string {
	h := md5.New()
	h.Write(bytestr.StringToBytes(val))
	return string(h.Sum(nil))
}
