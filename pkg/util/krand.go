package util

import (
	"math/rand"
)

const (
	KindNum   = "0123456789"
	KindLower = "abcdefghijklmnopqrstuvwxyz"
	KindUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	KindAll   = KindNum + KindLower + KindUpper
)

// Krand 生成指定长度和类型的随机字符串
func Krand(size int, kind string) string {
	if size <= 0 {
		return ""
	}

	// 根据传入的 kind 选择对应的字符集
	charset := kind
	if charset == "" {
		charset = KindAll
	}

	result := make([]byte, size)
	charsetLen := len(charset)

	// 只需要一次函数调用即可获取随机字节
	for i := range result {
		result[i] = charset[rand.Intn(charsetLen)]
	}

	return string(result)
}
