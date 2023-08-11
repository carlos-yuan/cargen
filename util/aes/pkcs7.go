/**
* @Author: qinwei
* @Date:   2023/3/3 09:07
* @FileUrl:   pkcs7.go
* @Desc:   PKCS7填充模式
 */

package aes

import (
	"fmt"
)

// PKCS7 填充
func pkcs7Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(plaintext, padText...)
}

// PKCS7 去除填充
func pkcs7UnPadding(padded []byte) ([]byte, error) {
	length := len(padded)
	unpadding := int(padded[length-1])
	if unpadding > length {
		return nil, fmt.Errorf("invalid padding")
	}
	return padded[:length-unpadding], nil
}
