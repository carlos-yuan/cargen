package ctl

import "github.com/carlos-yuan/cargen/util/config"

type Token interface {
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
