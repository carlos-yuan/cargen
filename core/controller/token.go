package ctl

import "github.com/carlos-yuan/cargen/core/config"

type Token interface {
	Clone() Token
	Sign() error
	Verify(token string) error
	GetPayLoad() Payload
	GetConfig() config.Token
}

type Payload interface {
	ToBytes() []byte
	FromBytes(from []byte) error
	Expire() int64
}
