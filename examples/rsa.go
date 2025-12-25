package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"

	"github.com/zngw/golib/crypt"
)

// GenerateRSAKeys 生成 RSA 密钥对
func GenerateRSAKeys() (privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, err error) {
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	publicKey = &privateKey.PublicKey
	return privateKey, publicKey, nil
}

// PrivateKeyToPEM 将私钥转换为 PEM 格式
func PrivateKeyToPEM(privateKey *rsa.PrivateKey) string {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return string(privateKeyPEM)
}

// PublicKeyToPEM 将公钥转换为 PEM 格式
func PublicKeyToPEM(publicKey *rsa.PublicKey) string {
	publicKeyBytes, _ := x509.MarshalPKIXPublicKey(publicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return string(publicKeyPEM)
}

func main() {
	// 1. 生成密钥对
	privateKey, publicKey, err := GenerateRSAKeys()
	if err != nil {
		log.Fatal("生成密钥失败:", err)
	}

	// 2. 转换为 PEM 格式
	privatePEM := PrivateKeyToPEM(privateKey)
	publicPEM := PublicKeyToPEM(publicKey)

	fmt.Println("私钥 (请妥善保管):")
	fmt.Println(privatePEM)
	fmt.Println("\n公钥 (可分发给前端):")
	fmt.Println(publicPEM)

	// 3. 模拟加密解密流程
	originalMessage := "mySecretPassword123"
	fmt.Printf("\n原始消息: %s\n", originalMessage)

	// 前端用公钥加密
	encrypted, err := crypt.RSAEncryptWithPublicKey(originalMessage, publicPEM)
	if err != nil {
		log.Fatal("加密失败:", err)
	}
	fmt.Printf("加密后 (Base64): %s\n", encrypted)

	// 后端用私钥解密
	decrypted, err := crypt.RSADecryptWithPrivateKey(encrypted, privatePEM)
	if err != nil {
		log.Fatal("解密失败:", err)
	}
	fmt.Printf("解密后: %s\n", decrypted)

	// 验证
	if originalMessage == decrypted {
		fmt.Println("\n✅ RSA 加解密成功！")
	} else {
		fmt.Println("\n❌ RSA 加解密失败！")
	}
}
