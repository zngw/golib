package crypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
)

// LoadPrivateKey 从字符串中加载RSA私钥
func LoadPrivateKey(prvKey string) (*rsa.PrivateKey, error) {
	// 检查是否包含-----BEGIN RSA PRIVATE KEY-----
	if strings.Index(prvKey, "-----BEGIN RSA PRIVATE KEY-----") < 0 {
		prvKey = fmt.Sprintf("-----BEGIN RSA PRIVATE KEY-----\n%v\n-----END PUBLIC KEY-----", prvKey)
	}
	block, _ := pem.Decode([]byte(prvKey))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// LoadPublicKey 从字符串中加载RSA公钥
func LoadPublicKey(pubKey string) (*rsa.PublicKey, error) {
	// 检查是否包含-----BEGIN PUBLIC KEY-----
	if strings.Index(pubKey, "-----BEGIN PUBLIC KEY-----") < 0 {
		pubKey = fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%v\n-----END PUBLIC KEY-----", pubKey)
	}

	block, _ := pem.Decode([]byte(pubKey))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey.(*rsa.PublicKey), nil
}

// RSAEncryptWithPublicKey 用公钥加密（OAEP填充）
func RSAEncryptWithPublicKey(message string, pubKey string) (string, error) {
	publicKey, err := LoadPublicKey(pubKey)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	encryptedBytes, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, []byte(message), nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

// RSADecryptWithPrivateKey 用私钥解密（OAEP填充）
func RSADecryptWithPrivateKey(encryptedMessage string, prvKey string) (string, error) {
	privateKey, err := LoadPrivateKey(prvKey)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	decryptedBytes, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, encryptedBytes, nil)
	if err != nil {
		return "", err
	}

	return string(decryptedBytes), nil
}

func SHA256WithRSASign(content, prvKey string) (sign string, err error) {
	privateKey, err := LoadPrivateKey(prvKey)
	if err != nil {
		return
	}

	h := sha256.New()
	h.Write([]byte(content))
	hashed := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return
	}
	sign = base64.StdEncoding.EncodeToString(signature)
	return
}

func SHA256WithRSAVerify(origdata, ciphertext, publicKey string) (bool, error) {
	pub, err := LoadPublicKey(publicKey)
	if err != nil {
		return false, err
	}

	body, err := base64.StdEncoding.DecodeString(ciphertext)

	h := sha256.New()
	h.Write([]byte(origdata))
	digest := h.Sum(nil)

	if err != nil {
		return false, err
	}
	err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest, body)
	if err != nil {
		return false, err
	}
	return true, nil
}
