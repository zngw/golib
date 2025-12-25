package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// RSAEncryptWithPublicKey 用公钥加密
func RSAEncryptWithPublicKey(message string, publicKeyPEM string) (string, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM block")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), []byte(message))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

// RSADecryptWithPrivateKey 用私钥解密
func RSADecryptWithPrivateKey(encryptedMessage string, privateKeyPEM string) (string, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedBytes)
	if err != nil {
		return "", err
	}

	return string(decryptedBytes), nil
}
