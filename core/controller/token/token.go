package token

import (
	e "comm/error"
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
)

var (
	ErrTokenMalformed        = errors.New("token is malformed")
	ErrTokenUnverifiable     = errors.New("token is unverifiable")
	ErrTokenSignatureInvalid = errors.New("token signature is invalid")

	ErrTokenInvalidAudience  = errors.New("token has invalid audience")
	ErrTokenExpired          = errors.New("token is expired")
	ErrTokenUsedBeforeIssued = errors.New("token used before issued")
	ErrTokenInvalidIssuer    = errors.New("token has invalid issuer")
	ErrTokenNotValidYet      = errors.New("token is not valid yet")
	ErrTokenInvalidId        = errors.New("token has invalid id")
	ErrTokenInvalidClaims    = errors.New("token has invalid claims")

	ErrTokenMethod = errors.New("token method invalid")
)

type Token struct {
	Type    Type
	Payload *Payload
	Token   string
}

func (t *Token) Sign(key interface{}) error {
	switch t.Type.Alg {
	case "HS256":
		return t.SignHS256(key)
	}
	return ErrTokenMethod
}

func (t *Token) Verify(key interface{}) error {
	switch t.Type.Alg {
	case "HS256":
		return t.verifyHS256(key)
	}
	return ErrTokenMethod
}

func (t *Token) GetPayLoad() *Payload {
	if t.Payload == nil {
		t.Payload = &Payload{}
		switch t.Type.Alg {
		case "HS256":
			b, err := DecodeSegment(t.Token[strings.Index(t.Token, ".")+1 : strings.LastIndex(t.Token, ".")])
			if err != nil {
				log.Fatal(e.GetSite(1), err)
			}
			err = json.Unmarshal(b, t.Payload)
			if err != nil {
				log.Fatal(e.GetSite(1), err)
			}
		default:
			panic(ErrTokenMethod)
		}
	}
	return t.Payload
}

func (t *Token) SignHS256(key interface{}) error {
	typeStr := EncodeSegment(t.Type.JSON())
	payloadStr := EncodeSegment(t.Payload.JSON())
	b := make([]byte, 0, len(typeStr)+len(payloadStr)+45)
	b = append(b, typeStr...)
	b = append(b, byte('.'))
	b = append(b, payloadStr...)
	sign, err := HS256Sign(b, []byte(key.(string)))
	if err != nil {
		return err
	}
	b = append(b, byte('.'))
	b = append(b, sign...)
	t.Token = string(b)
	return err
}

func (t *Token) verifyHS256(key interface{}) error {
	signIndex := strings.LastIndex(t.Token, ".")
	if signIndex == -1 {
		return ErrTokenNotValidYet
	}
	signStr := t.Token[signIndex+1:]
	signTypePayload := t.Token[:signIndex]
	err := HS256Verify(signTypePayload, signStr, []byte(key.(string)))
	return err
}

type Type struct {
	Alg interface{} `json:"alg"`
	Typ string      `json:"typ"`
}

func (t *Type) JSON() []byte {
	b, _ := json.Marshal(t)
	return b
}

type Payload struct {
	Datascope    string `json:"datascope,omitempty"`
	EnterpriseId int64  `json:"enterpriseId,string,omitempty"`
	Exp          int64  `json:"exp,omitempty"`
	GrpcToken    string `json:"grpcToken,omitempty"`
	Identity     int    `json:"identity,omitempty"`
	Nice         string `json:"nice,omitempty"`
	OrigIat      int    `json:"orig_iat,omitempty"`
	OpId         int64  `json:"opId,string,omitempty"` //获取时 使用GetOpId 避免鉴权出错
	Role         int32  `json:"role,omitempty"`
	Roleid       int    `json:"roleid,omitempty"`
	Rolekey      string `json:"rolekey,omitempty"`
	Rolename     string `json:"rolename,omitempty"`
	UserId       int64  `json:"userId,string,omitempty"`
}

func (p *Payload) JSON() []byte {
	b, _ := json.Marshal(p)
	return b
}

// 绑定许可证编号
func (p *Payload) BindOpId(id *int64) {
	if p.OpId == 0 {
		panic(e.ParamsDealError.SetErr(nil, "许可证信息尚未选择"))
	}
	*id = p.OpId
}

// 带鉴权获得许可证编号
func (p *Payload) GetOpId() int64 {
	if p.OpId == 0 {
		panic(e.ParamsDealError.SetErr(nil, "许可证信息尚未选择"))
	}
	return p.OpId
}

var GrpcTokenKey struct{}

const GrpcTokenStringKey = "grpc_token"

func (p Payload) GrpcTokenContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, GrpcTokenKey, p.GrpcToken)
}
