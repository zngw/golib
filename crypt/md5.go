package crypt

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/zngw/golib/str"
)

func Md5Hex(plaintext any) string {
	m := md5.New()
	m.Write(str.ToBytes(plaintext))
	return hex.EncodeToString(m.Sum(nil))
}
