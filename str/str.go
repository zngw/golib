// Package str
// @Description 一些字符串与数字转换方法
// @Author  55
// @Date  2022/5/30
package str

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

// StringToInt 字符串转整型
func StringToInt(str string) (i int) {
	i64, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return
	}

	i = int(i64)
	return
}

// StringToFloat 字符串转浮点型
func StringToFloat(str string) (f float64) {
	f, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return
	}

	return
}

// ToString any类型转成字符串，非基础类型用json格式
func ToString(a any) string {
	switch v := a.(type) {
	case bool:
		if v {
			return "true"
		} else {
			return "false"
		}
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uintptr:
		return strconv.FormatUint(uint64(v), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case complex64, complex128:
		return fmt.Sprintf("%v", v)
	case []byte:
		return string(v)
	default:
		// 非基础类型转成 json
		d, _ := json.Marshal(v)
		return string(d)
	}
}

// ToBytes []byte直接返回，其他的先转成字符串再转[]byte
func ToBytes(a any) []byte {
	switch v := a.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	default:
		return []byte(ToString(a))
	}
}

var seedCounter int64

// RandString 生成指定长度的随机字符串
func RandString(codeLen int) string {
	// 1. 定义原始字符串
	rawStr := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	// 2. 定义一个buf，并且将buf交给bytes往buf中写数据
	buf := make([]byte, 0, codeLen)
	b := bytes.NewBuffer(buf)
	// 随机从中获取
	random := rand.New(rand.NewSource(time.Now().UnixNano() ^ int64(os.Getegid()) ^ atomic.AddInt64(&seedCounter, 1)))
	for rawStrLen := len(rawStr); codeLen > 0; codeLen-- {
		randNum := random.Intn(rawStrLen)
		b.WriteByte(rawStr[randNum])
	}
	return b.String()
}
