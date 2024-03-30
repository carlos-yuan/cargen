package gmsm

import (
	"encoding/hex"
	ZZMarquis_sm2 "github.com/ZZMarquis/gm/sm2"
	"github.com/deatil/go-cryptobin/cryptobin/sm2"
)

// 数据加密
func Sm2Encrypt(publicKey, data string) (string, error) {
	pk, err := hex.DecodeString(publicKey)
	if err != nil {
		return "", err
	}
	pub := sm2.NewSM2().FromPublicKeyBytes(pk).GetPublicKey()
	pbByte := make([]byte, 0, 64)
	pbByte = append(pbByte, pub.X.Bytes()...)
	pbByte = append(pbByte, pub.Y.Bytes()...)
	pbk, err := ZZMarquis_sm2.RawBytesToPublicKey(pbByte)
	if err != nil {
		return "", err
	}
	rb, err := ZZMarquis_sm2.Encrypt(pbk, []byte(data), ZZMarquis_sm2.C1C3C2)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(rb), err
}

// 数据解密
func Sm2Decrypt(sk, data string) (string, error) {
	pk, err := hex.DecodeString(sk)
	if err != nil {
		return "", err
	}
	pri := sm2.NewSM2().FromPrivateKeyBytes(pk).GetPrivateKey()
	priByte := make([]byte, 0, 64)
	priByte = append(priByte, pri.D.Bytes()...)
	prik, err := ZZMarquis_sm2.RawBytesToPrivateKey(pri.D.Bytes())
	if err != nil {
		return "", err
	}
	b, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	rb, err := ZZMarquis_sm2.Decrypt(prik, b, ZZMarquis_sm2.C1C3C2)
	if err != nil {
		return "", err
	}
	return string(rb), err
}
