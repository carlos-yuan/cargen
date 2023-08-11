package rd

import (
	"comm/timeUtil"
	"crypto/rand"
	uuid "github.com/satori/go.uuid"
	"math/big"
	mrand "math/rand"
	"strconv"
	"time"
)

func GetUUID() string {
	return uuid.NewV4().String()
}

func RandInt64(min, max int64) int64 {
	maxBigInt := big.NewInt(max)
	i, _ := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < min {
		RandInt64(min, max)
	}
	return i.Int64()
}

func GetIntNumber(l int) int64 {
	if l > 20 {
		panic("GetIntNumber err")
	}
	num, _ := strconv.ParseInt(GetNumber(l), 10, 64)
	return num
}

func GetNumber(l int) string {
	str := "123456789"
	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	return GetOrderId(string(str[r.Intn(len(str))]), l)
}

func GetOrderId(head string, l int) string {
	str := "0123456789"
	result := []byte(head)
	var temp = GetUUID()
	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l-len(head); i++ {
		t := temp[r.Intn(len(temp)-1)] % 10
		result = append(result, str[t])
	}
	return string(result)
}

func GetDateId(l int) string {
	return GetOrderId(timeUtil.ToyyyyMMddHHmmss(), l)
}

func GetDateIdInt(l int) int64 {
	str := GetDateId(l)
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}
