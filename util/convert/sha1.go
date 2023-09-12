package convert

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
)

func Sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return hex.EncodeToString(t.Sum(nil))
}
