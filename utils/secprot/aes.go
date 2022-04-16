package secprot

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
	0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

var (
	aessalt []byte
	cfbEnc  cipher.Stream
	cfbDec  cipher.Stream
)

func SetAESSalt(salt string) {
	aessalt = []byte(Md5Encode(salt))
}

func AesCfbEnc(origin []byte) string {
	block, err := aes.NewCipher(aessalt)
	if err != nil {
		panic(err)
	}
	enc := make([]byte, aes.BlockSize+len(origin))
	iv := enc[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(enc[aes.BlockSize:], origin)
	return Base64Encode(string(enc))
}

func AesCfbDec(enc string) string {
	vb, err := Base64Decode(enc)
	if err != nil {
		return ""
	}

	block, _ := aes.NewCipher(aessalt)
	if len(enc) < aes.BlockSize {
		return ""
	}
	iv := vb[:aes.BlockSize]
	vb = vb[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(vb, vb)
	return string(vb)
}
