package convert

import (
	"bytes"
	"strconv"
	"strings"
	"unsafe"
)

func ParseInt(number string) int64 {
	i, _ := strconv.ParseInt(number, 10, 64)
	return i
}

// 限制字符串长度
func StrLimit(str, suffix string, length int32) string {
	if len(str) > int(length) {
		str = str[0:int(length)-len(suffix)] + suffix
	}
	return str
}

// 整型转固定长度字符串
func IntToCode(code int, length int32) string {
	var buf bytes.Buffer
	codeStr := strconv.Itoa(code)
	for i := 0; i < int(length)-len(codeStr); i++ {
		buf.WriteString("0")
	}
	buf.WriteString(strconv.Itoa(code))
	return buf.String()
}

func FormatCode(code, separate string) string {
	if len(code) == 3 {
		return code
	} else if len(code) == 7 {
		return code[:3] + separate + code[3:]
	} else if len(code) == 12 {
		return code[:3] + separate + code[3:7] + separate + code[7:]
	} else if len(code) == 18 {
		return code[:3] + separate + code[3:7] + separate + code[7:12] + separate + code[12:18]
	}
	return ""
}

func SplitCode(code, separate string) []string {
	var res []string
	if separate != "" {
		return strings.Split(code, separate)
	} else {
		if len(code) == 3 {
			res = append(res, code)
		} else if len(code) == 7 {
			res = append(res, code[:3], code[3:])
		} else if len(code) == 12 {
			res = append(res, code[:3], code[3:7], code[7:])
		} else if len(code) == 18 {
			res = append(res, code[:3], code[3:7], code[7:12], code[12:18])
		}
	}
	return res
}

func SplitAreaCode(code string) []string {
	var res []string
	if len(code) == 3 {
		res = append(res, code)
	} else if len(code) == 7 {
		res = append(res, code[:3], code[:7])
	} else if len(code) == 12 {
		res = append(res, code[:3], code[:7], code[:12])
	} else if len(code) == 18 {
		res = append(res, code[:3], code[:7], code[:12], code[:18])
	}
	return res
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func RemoveSpecialCharacters(str string) string {
	s := []rune(str)
	for i := range s {
		if s[i] > 40959 {
			s[i] = ' '
		}
	}
	return string(s)
}

func ToI32(numStr string) int32 {
	n, err := strconv.ParseInt(numStr, 10, 32)
	if err != nil {
		panic(err)
	}
	return int32(n)
}

func ToI64(numStr string) int64 {
	n, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		panic(err)
	}
	return int64(n)
}
