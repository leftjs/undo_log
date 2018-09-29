package util

import (
	"encoding/base64"
	"math"
	"math/rand"
	"strings"
)

const SALT_SEP = "#abasdf#" // 分隔符

func Encode(data []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(data))
}

func Decode(data []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(data))
}

// 带盐加密
func EncodeWithRandomSalt(data []byte) []byte {
	salt := []byte(string(rand.Intn(math.MaxInt32)))
	var pre []byte
	pre = append(pre, data...)
	pre = append(pre, []byte(SALT_SEP)...)
	pre = append(pre, salt...)
	return Encode(pre)
}

// 解密去盐
func DecodeFromRandomSalt(data []byte) ([]byte, error) {
	post, err := Decode(data)
	return []byte(strings.Split(string(post), SALT_SEP)[0]), err
}
