package main

import (
	"github.com/zngw/golib/crypt"
	"github.com/zngw/golib/log"
)

func main() {
	prvStr, pubStr, err := crypt.Ed25519GenerateKeyPair()
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info("generate keypair \n%s\n\n%s", prvStr, pubStr)

	message := "Hello Ed25519"

	// 签名
	signature, err := crypt.Ed25519Sign(message, prvStr)
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info("signature string:", signature)

	// 验证签
	ok, err := crypt.Ed25519Verify(message, signature, pubStr)
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info("verify result: %v", ok)
}
