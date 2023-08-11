package md5

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

func Encode(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func EncodeAny(data interface{}) string {
	h := md5.New()
	b, _ := json.Marshal(data)
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

func DoubleEncode(str string) string {
	return Encode(Encode(str))
}

func EncodeByte(bytes []byte) string {
	h := md5.New()
	h.Write(bytes)
	return hex.EncodeToString(h.Sum(nil))
}
