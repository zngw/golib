package crypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/zngw/golib/str"
)

func HmacSha256Hex(plaintext any, secret string) string {
	// 将密钥和消息转换为字节
	key := []byte(secret)
	data := str.ToBytes(plaintext)

	// 创建HMAC-SHA256哈希器
	h := hmac.New(sha256.New, key)

	// 写入消息数据
	h.Write(data)

	// 计算HMAC摘要
	mac := h.Sum(nil)

	// 将二进制结果转换为十六进制字符串
	return hex.EncodeToString(mac)
}
