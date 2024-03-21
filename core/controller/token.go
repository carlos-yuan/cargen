package ctl

import (
	"github.com/carlos-yuan/cargen/core/config"
	"io"
)

type Token interface {
	Clone() Token
	Sign() error
	Verify(token string) error                         //验证
	Decrypt(body io.ReadCloser) (io.ReadCloser, error) //解密
	GetPayLoad() Payload
	GetConfig() config.Token
}

type Payload interface {
	Clone() Payload
	ToBytes() []byte
	FromBytes(from []byte) error
	Expire() int64
}
