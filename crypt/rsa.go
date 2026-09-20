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
)

// LoadPrivateKey 从字符串中加载RSA私钥
func LoadPrivateKey(prvKey string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(prvKey))
	if block == nil {
		return nil, fmt.Errorf("解析 PEM 失败")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		// PKCS#1 格式
		return x509.ParsePKCS1PrivateKey(block.Bytes)

	case "PRIVATE KEY":
		// PKCS#8 格式
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("不是 RSA 私钥")
		}
		return rsaKey, nil

	default:
		return nil, fmt.Errorf("不支持的私钥类型: %s", block.Type)
	}
}

// LoadPublicKey 从字符串中加载RSA公钥
func LoadPublicKey(pubKey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubKey))
	if block == nil {
		return nil, fmt.Errorf("解析 PEM 失败")
	}

	switch block.Type {
	case "PUBLIC KEY":
		// PKCS#8 / SubjectPublicKeyInfo 格式
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("解析 PKIX 公钥失败: %v", err)
		}

		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("不是 RSA 公钥")
		}
		return rsaPub, nil

	case "RSA PUBLIC KEY":
		// PKCS#1 格式
		return x509.ParsePKCS1PublicKey(block.Bytes)

	default:
		return nil, fmt.Errorf("不支持的公钥类型: %s", block.Type)
	}
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
