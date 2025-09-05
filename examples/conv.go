package main

import (
	"github.com/zngw/golib/crypt"
	"github.com/zngw/golib/log"
)

func main() {
	key := "1234567890abcdef1234567890abcdef"
	text := "https://zengwu.com.cn"

	gcmCipher := crypt.GcmEncrypt(text, key)
	gcmPlaintext := crypt.GcmDecrypt(gcmCipher, key)
	log.Trace("gcm加密/解密：加密base64=%s，解密明文=%s", gcmCipher, gcmPlaintext)

	aesCipher := crypt.AesEcbEncrypt(text, key[:16])
	aesPlaintext := crypt.AesEcbDecrypt(aesCipher, key[:16])
	log.Trace("aes ecb加密/解密：加密base64=%s，解密明文=%s", aesCipher, aesPlaintext)
}
