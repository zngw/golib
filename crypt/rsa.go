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

// LoadPublicKey 加载公钥
// 支持：
// - PEM：-----BEGIN PUBLIC KEY-----
// - PEM：-----BEGIN RSA PUBLIC KEY-----
// - 纯 Base64 字符串
// - DER 字节
func LoadPublicKey(pubKey string) (*rsa.PublicKey, error) {
	if len(pubKey) == 0 {
		return nil, fmt.Errorf("公钥内容为空")
	}

	// 1. 先尝试 PEM 解析
	if block, _ := pem.Decode([]byte(pubKey)); block != nil {
		return parsePublicKeyPEM(block)
	}

	// 2. PEM 失败，按 Base64 / DER 解析
	return parsePublicKeyBase64OrDER([]byte(pubKey))
}

// parsePublicKeyPEM 根据 PEM 头解析公钥
func parsePublicKeyPEM(block *pem.Block) (*rsa.PublicKey, error) {
	switch block.Type {
	case "PUBLIC KEY":
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		return key.(*rsa.PublicKey), nil

	case "RSA PUBLIC KEY":
		return x509.ParsePKCS1PublicKey(block.Bytes)

	default:
		return nil, fmt.Errorf("不支持的公钥 PEM 类型: %s", block.Type)
	}
}

// parsePublicKeyBase64OrDER 处理没有 PEM 头的公钥
func parsePublicKeyBase64OrDER(data []byte) (*rsa.PublicKey, error) {
	der, err := decodeToDER(data)
	if err != nil {
		return nil, err
	}

	// 先尝试 PKIX 公钥
	if pub, err := x509.ParsePKIXPublicKey(der); err == nil {
		return pub.(*rsa.PublicKey), nil
	}

	// 再尝试 PKCS#1 RSA 公钥
	if pub, err := x509.ParsePKCS1PublicKey(der); err == nil {
		return pub, nil
	}

	return nil, fmt.Errorf("无法识别公钥格式")
}

// LoadPrivateKey 加载私钥
// 支持：
// - PEM：-----BEGIN PRIVATE KEY-----
// - PEM：-----BEGIN RSA PRIVATE KEY-----
// - 纯 Base64 字符串
// - DER 字节
func LoadPrivateKey(prvKey string) (*rsa.PrivateKey, error) {
	if len(prvKey) == 0 {
		return nil, fmt.Errorf("私钥内容为空")
	}

	// 1. 先尝试 PEM 解析
	if block, _ := pem.Decode([]byte(prvKey)); block != nil {
		return parsePrivateKeyPEM(block)
	}

	// 2. PEM 失败，按 Base64 / DER 解析
	return parsePrivateKeyBase64OrDER([]byte(prvKey))
}

// parsePrivateKeyPEM 根据 PEM 头解析私钥
func parsePrivateKeyPEM(block *pem.Block) (*rsa.PrivateKey, error) {
	switch block.Type {
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return key.(*rsa.PrivateKey), nil

	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)

	default:
		return nil, fmt.Errorf("不支持的私钥 PEM 类型: %s", block.Type)
	}
}

// parsePrivateKeyBase64OrDER 处理没有 PEM 头的私钥
func parsePrivateKeyBase64OrDER(data []byte) (*rsa.PrivateKey, error) {
	der, err := decodeToDER(data)
	if err != nil {
		return nil, err
	}

	// 依次尝试常见私钥格式
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		return key.(*rsa.PrivateKey), nil
	}

	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}

	return nil, fmt.Errorf("无法识别私钥格式")
}

// decodeToDER 把输入统一转成 DER 字节
// 优先按 Base64 解码；如果不是合法 Base64，则当作原始 DER
func decodeToDER(data []byte) ([]byte, error) {
	trimmed := strings.TrimSpace(string(data))

	der, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		// 不是 Base64，就当作原始 DER
		return data, nil
	}

	return der, nil
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
