package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/zngw/golib/str"
)

// GcmEncrypt 使用Aes-Gcm加密明文，出错时返回空字符串。
// 加密成功时输出 base64(初始化向量+密文)
func GcmEncrypt(plaintext, secretKey string) (string, error) {
	// 需要解码
	key, err := hex.DecodeString(secretKey)
	if err != nil {
		return "", fmt.Errorf("invalid secret key hex: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, str.ToBytes(plaintext), nil)

	encoded := base64.StdEncoding.EncodeToString(cipherText)
	return encoded, nil
}

// GcmDecrypt 使用Aes-Gcm解密，出错时返回空字符串和错误
// ciphertext 为 base64(初始化向量+密文)
func GcmDecrypt(ciphertext, secretKey string) (string, error) {
	key, err := hex.DecodeString(secretKey)
	if err != nil {
		return "", fmt.Errorf("invalid secret key hex: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	encryptBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(encryptBytes) < gcm.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}

	plaintext, err := gcm.Open(nil, encryptBytes[:gcm.NonceSize()], encryptBytes[gcm.NonceSize():], nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// AesEcbEncrypt [遗留兼容] ECB模式不安全，相同明文块产生相同密文块，新代码请使用 GcmEncrypt。
// 使用pkcs填充，以base64格式输出
func AesEcbEncrypt(plaintext, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	data := pkcs5Padding([]byte(plaintext), block.BlockSize())
	decrypted := make([]byte, len(data))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(decrypted[bs:be], data[bs:be])
	}

	return base64.StdEncoding.EncodeToString(decrypted), nil
}

// AesEcbDecrypt [遗留兼容] ECB模式不安全，新代码请使用 GcmDecrypt。
// 使用pkcs填充，秘文为base64格式输入
func AesEcbDecrypt(ciphertext, key string) (result string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("aes ecb decrypt panic: %v", r)
		}
	}()

	encryptBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return
	}

	decrypted := make([]byte, len(encryptBytes))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(encryptBytes); bs, be = bs+size, be+size {
		block.Decrypt(decrypted[bs:be], encryptBytes[bs:be])
	}

	return string(pkcs5UnPadding(decrypted)), nil
}

// // pkcs5补码算法
func pkcs5UnPadding(origData []byte) []byte {
	// 1. 计算数据的总长度
	length := len(origData)
	if length == 0 {
		return origData
	}
	// 2. 根据填充的字节值得到填充的次数
	number := int(origData[length-1])
	if number > length {
		return origData
	}
	// 3. 将尾部填充的number个字节去掉
	return origData[:(length - number)]
}

// pks5填充的尾部数据
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
