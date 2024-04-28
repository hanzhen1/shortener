package base62

import (
	"math"
	"strings"
)

// 62进制转换模块
// 0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
// 0-9:0-9
// a-z:10-35
// A-Z:36-61
// 10进制数 	转换   62进制数
//
//	1                 1
//	9                 9
//	10                a
//	11                b
//	61                Z
//	62                10
//	63                11
//	6347              ? 1En
//
// 如何实现62进制转换
// const base62Str = `0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
//const base62Str = `aWc9de0fgSi2jklNmn1opqEhUsbTvw3yzABC4DFGHxIJrKuL5MOP6QtR7VXY8Z`

// 为了避免被人恶意请求，我们可以将上面的字符串打乱
var base62Str string

// MustInit 要使用base62这个包必须要调用该函数完成初始化
func MustInit(bs string) {
	if len(bs) == 0 {
		panic("need base string")
	}
	base62Str = bs
}

// IntToString 十进制数转为62进制字符串
func IntToString(seq uint64) string {
	if seq == 0 {
		return string(base62Str[0])
	}
	slice := []byte{}
	for seq > 0 { // 6347-> 23 40 1
		mod := seq % 62
		div := seq / 62
		slice = append(slice, base62Str[mod])
		seq = div
	}
	//最后把得到的数反转一下
	return string(reverse(slice))
}

// StringToInt 62进制字符串转为十进制数
func StringToInt(s string) (req uint64) {
	slice := []byte(s)
	slice = reverse(slice) //反转后 nE1
	for idx, b := range slice {
		base := math.Pow(62, float64(idx))
		req += uint64(strings.Index(base62Str, string(b))) * uint64(base)
	}
	return req
}

// 反转切片
func reverse(s []byte) []byte {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
