package gmsm

import (
	"encoding/hex"
	"github.com/tjfoc/gmsm/sm3"
	"strings"
)

func SM3EncodeString(data []byte) string {
	return strings.ToUpper(hex.EncodeToString(sm3.Sm3Sum(data)))
}
