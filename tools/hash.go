package tools

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

//md5加密
func Md5(src string) string {
	m := md5.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

//Sha256加密
func Sha256(src string) string {
	m := sha256.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}
// Base64Encode 标准 base64 编码
func Base64Encode(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}

func Base64Decode(str string) string {
	// 兼容带填充（StdEncoding）和不带填充（RawStdEncoding）的 base64
	if data, err := base64.StdEncoding.DecodeString(str); err == nil {
		return string(data)
	}
	if data, err := base64.RawStdEncoding.DecodeString(str); err == nil {
		return string(data)
	}
	return ""
}
