package crypt

import (
	"crypto/cipher"
	"crypto/des"
)

// DesEcbEncrypt [遗留兼容] DES密钥仅56位有效长度，安全性不足，新代码请使用 AES-GCM。
func DesEcbEncrypt(data, key []byte) []byte {
	//NewCipher创建一个新的加密块
	block, err := des.NewCipher(key)
	if err != nil {
		return nil
	}

	bs := block.BlockSize()
	data = pkcs5Padding(data, bs)
	if len(data)%bs != 0 {
		return nil
	}

	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		//Encrypt加密第一个块，将其结果保存到dst
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out
}

// DesEcbDecrypt [遗留兼容] DES密钥仅56位有效长度，安全性不足，新代码请使用 AES-GCM。
func DesEcbDecrypt(data, key []byte) []byte {
	//NewCipher创建一个新的加密块
	block, err := des.NewCipher(key)
	if err != nil {
		return nil
	}

	bs := block.BlockSize()
	if len(data)%bs != 0 {
		return nil
	}

	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		//Encrypt加密第一个块，将其结果保存到dst
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}

	out = pkcs5UnPadding(out)

	return out
}

// DesCbcEncrypt [遗留兼容] DES密钥仅56位有效长度，安全性不足，新代码请使用 AES-GCM。
func DesCbcEncrypt(data, key, iv []byte) []byte {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil
	}

	data = pkcs5Padding(data, block.BlockSize())
	cryptText := make([]byte, len(data))

	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(cryptText, data)
	return cryptText
}

// DesCbcDecrypt [遗留兼容] DES密钥仅56位有效长度，安全性不足，新代码请使用 AES-GCM。
func DesCbcDecrypt(data, key, iv []byte) []byte {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	cryptText := make([]byte, len(data))
	blockMode.CryptBlocks(cryptText, data)
	cryptText = pkcs5UnPadding(cryptText)

	return cryptText
}

// DesCtrEncrypt [遗留兼容] DES密钥仅56位有效长度，安全性不足，新代码请使用 AES-GCM。
func DesCtrEncrypt(data, key, iv []byte) []byte {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil
	}

	data = pkcs5Padding(data, block.BlockSize())
	cryptText := make([]byte, len(data))

	blockMode := cipher.NewCTR(block, iv)
	blockMode.XORKeyStream(cryptText, data)
	return cryptText
}

// DesCtrDecrypt [遗留兼容] DES密钥仅56位有效长度，安全性不足，新代码请使用 AES-GCM。
func DesCtrDecrypt(data, key, iv []byte) []byte {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil
	}

	blockMode := cipher.NewCTR(block, iv)
	cryptText := make([]byte, len(data))
	blockMode.XORKeyStream(cryptText, data)
	cryptText = pkcs5UnPadding(cryptText)

	return cryptText
}

// DesOfbEncrypt [遗留兼容] DES密钥仅56位有效长度，安全性不足，新代码请使用 AES-GCM。
func DesOfbEncrypt(data, key, iv []byte) []byte {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil
	}

	data = pkcs5Padding(data, block.BlockSize())
	cryptText := make([]byte, len(data))

	blockMode := cipher.NewOFB(block, iv)
	blockMode.XORKeyStream(cryptText, data)
	return cryptText
}

// DesOfbDecrypt [遗留兼容] DES密钥仅56位有效长度，安全性不足，新代码请使用 AES-GCM。
func DesOfbDecrypt(data, key, iv []byte) []byte {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil
	}

	blockMode := cipher.NewOFB(block, iv)
	cryptText := make([]byte, len(data))
	blockMode.XORKeyStream(cryptText, data)
	cryptText = pkcs5UnPadding(cryptText)

	return cryptText
}

// DesCfbEncrypt [遗留兼容] DES密钥仅56位有效长度，安全性不足，新代码请使用 AES-GCM。
func DesCfbEncrypt(data, key, iv []byte) []byte {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil
	}

	data = pkcs5Padding(data, block.BlockSize())
	cryptText := make([]byte, len(data))

	blockMode := cipher.NewCFBEncrypter(block, iv)
	blockMode.XORKeyStream(cryptText, data)
	return cryptText
}

// DesCfbDecrypt [遗留兼容] DES密钥仅56位有效长度，安全性不足，新代码请使用 AES-GCM。
func DesCfbDecrypt(data, key, iv []byte) []byte {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil
	}

	blockMode := cipher.NewCFBDecrypter(block, iv)
	cryptText := make([]byte, len(data))
	blockMode.XORKeyStream(cryptText, data)
	cryptText = pkcs5UnPadding(cryptText)

	return cryptText
}
