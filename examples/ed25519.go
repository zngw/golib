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

	pub, err := crypt.Ed25519ParsePublicKey(pubStr)
	if err != nil {
		log.Error(err.Error())
		return
	}

	prv, err := crypt.Ed25519ParsePrivateKey(prvStr)
	if err != nil {
		log.Error(err.Error())
		return
	}

	message := "Hello Ed25519"

	// 签名
	signature, err := crypt.Ed25519Sign(prv, message)
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info("signature string:", signature)

	// 验证签
	ok, err := crypt.Ed25519Verify(pub, message, signature)
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info("verify result: %v", ok)
}
