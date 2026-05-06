package main

import (
	"github.com/zngw/golib/crypt"
	"github.com/zngw/golib/log"
)

func main() {
	key := "1234567890abcdef1234567890abcdef"
	text := "https://zengwu.com.cn"

	gcmCipher, err := crypt.GcmEncrypt(text, key)
	if err != nil {
		log.Error("gcm加密失败: %v", err)
		return
	}
	gcmPlaintext, err := crypt.GcmDecrypt(gcmCipher, key)
	if err != nil {
		log.Error("gcm解密失败: %v", err)
		return
	}
	log.Trace("gcm加密/解密：加密base64=%s，解密明文=%s", gcmCipher, gcmPlaintext)

	aesCipher, err := crypt.AesEcbEncrypt(text, key[:16])
	if err != nil {
		log.Error("aes ecb加密失败: %v", err)
		return
	}
	aesPlaintext, err := crypt.AesEcbDecrypt(aesCipher, key[:16])
	if err != nil {
		log.Error("aes ecb解密失败: %v", err)
		return
	}
	log.Trace("aes ecb加密/解密：加密base64=%s，解密明文=%s", aesCipher, aesPlaintext)
}
