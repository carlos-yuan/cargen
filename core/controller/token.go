package ctl

import (
	"github.com/carlos-yuan/cargen/core/config"
)

type Token interface {
	Clone() Token
	Sign() error
	Verify(ctx ControllerContext) error //验证
	GetPayLoad() Payload
	GetConfig() config.Token
}

type Payload interface {
	Clone() Payload
	ToBytes() []byte
	FromBytes(from []byte) error
	Expire() int64
}
