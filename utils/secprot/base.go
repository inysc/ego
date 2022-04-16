package secprot

import (
	"encoding/base32"
	"encoding/base64"

	"github.com/inysc/ego/utils/bytestr"
)

func Base32Encode(key string) string {
	return base32.StdEncoding.EncodeToString(bytestr.StringToBytes(key))
}

func Base64Encode(key string) string {
	return base64.StdEncoding.EncodeToString(bytestr.StringToBytes(key))
}

func Base32Decode(val string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(val)
}

func Base64Decode(val string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(val)
}
