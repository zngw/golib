package crypt

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

// Ed25519GenerateKeyPair 生成 Ed25519 密钥对
func Ed25519GenerateKeyPair() (prv, pub string, err error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return
	}

	// 私钥匙转 PKCS#8格式
	prvDer, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return
	}

	prvBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: prvDer,
	}

	// 公钥转PKIX 格式
	pubDer, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return
	}

	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubDer,
	}

	prv = string(pem.EncodeToMemory(prvBlock))
	pub = string(pem.EncodeToMemory(pubBlock))
	return
}

// Ed25519ParsePublicKey 从 hex 字符串解析 Ed25519 公钥
func Ed25519ParsePublicKey(publicKey string) (ed25519.PublicKey, error) {
	if publicKey == "" {
		return nil, errors.New("public key string is empty")
	}

	if strings.Index(publicKey, "BEGIN PUBLIC KEY") < 0 {
		publicKey = fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%v\n-----END PUBLIC KEY-----", publicKey)
	}

	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("failed to decode public key pem")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pub, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("parsed key is not ed25519 public key")
	}

	return pub, nil
}

// Ed25519ParsePrivateKey 从 hex 字符串解析 Ed25519 私钥
func Ed25519ParsePrivateKey(privateKey string) (ed25519.PrivateKey, error) {
	if privateKey == "" {
		return nil, errors.New("private key string is empty")
	}

	if strings.Index(privateKey, "END PRIVATE KEY") < 0 {
		privateKey = fmt.Sprintf("-----BEGIN PRIVATE KEY-----\n%v\n-----END PRIVATE KEY-----", privateKey)
	}

	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("failed to decode private key pem")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	priv, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("parsed key is not ed25519 private key")
	}

	return priv, nil
}

// Ed25519Sign 使用私钥对消息进行签名
func Ed25519Sign(privateKey ed25519.PrivateKey, message string) (string, error) {
	if len(privateKey) != ed25519.PrivateKeySize {
		return "", errors.New("invalid ed25519 private key")
	}

	if message == "" {
		return "", errors.New("message is nil")
	}

	signature := ed25519.Sign(privateKey, []byte(message))

	signStr := base64.StdEncoding.EncodeToString(signature)
	return signStr, nil
}

// Ed25519Verify 使用公钥验证签名
func Ed25519Verify(publicKey ed25519.PublicKey, message string, signature string) (bool, error) {
	if len(publicKey) != ed25519.PublicKeySize {
		return false, errors.New("invalid ed25519 public key")
	}

	if message == "" {
		return false, errors.New("message is nil")
	}

	body, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}

	if len(body) != ed25519.SignatureSize {
		return false, errors.New("invalid ed25519 signature")
	}

	return ed25519.Verify(publicKey, []byte(message), body), nil
}
