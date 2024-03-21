package gmsm

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/tjfoc/gmsm/sm4"
)

// PKCSPadd填充算法
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	fmt.Printf("PKCS7Padding blockSize[%d]\n", blockSize)
	fmt.Printf("PKCS7Padding len(ciphertext)[%d]\n", len(ciphertext))
	padding := blockSize - len(ciphertext)%blockSize
	fmt.Printf("PKCS7Padding len(ciphertext) blockSize[%d]\n", len(ciphertext)%blockSize)

	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	fmt.Println(padText)
	return append(ciphertext, padText...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	fmt.Println(origData)
	fmt.Printf("PKCSUnPadding origData len %d\n", len(origData))
	unPadding := int(origData[length-1])
	fmt.Printf("PKCSUnPadding origData unPadding %d\n", unPadding)
	return origData[:(length - unPadding)]
}

// 金融sm4加密ECB模式不填充，注意被加密数据必须为16个字节
func SM4ECBEncrypt(originalBytes, key []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		panic(err)
	}
	fmt.Println("SM44ECBEncrypt block.BlockSize")

	//originalBytes = PKCS7Padding(originalBytes, block.BlockSize())

	if len(originalBytes) != 16 {
		fmt.Println("长度不是16")
		panic(err)
		return nil, err
	}

	cipherArr := make([]byte, len(originalBytes))
	block.Encrypt(cipherArr, originalBytes)
	return cipherArr, nil

}

// 金融sm4解密ECB模式不填充，注意被加密数据必须为16个字节
func SM4ECBDecrypt(originalBytes, key []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(originalBytes) != 16 {
		fmt.Println("长度不是16")
		panic(err)
		return nil, err
	}

	cipherArr := make([]byte, len(originalBytes))
	block.Decrypt(cipherArr, originalBytes)
	//cipherArr = PKCS7UnPadding(cipherArr)
	return cipherArr, nil

}

func SM4ECBBase64Encrypt(originalText string, key []byte) (string, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		panic(err)
	}
	originalBytes := []byte(originalText)
	originalBytes = PKCS7Padding(originalBytes, block.BlockSize())

	cipherArr := make([]byte, 0)
	cArr := make([]byte, 16)

	j := 0
	for i := 0; i < len(originalBytes)/16; i++ {
		original := originalBytes[j : j+16]
		block.Encrypt(cArr, original)
		cipherArr = append(cipherArr, cArr...)
		j = j + 16
	}
	base64Str := base64.StdEncoding.EncodeToString(cipherArr)
	return base64Str, nil

}

func SM4ECBBase64Decrypt(cipherText string, key []byte) (string, error) {
	cInArr, _ := base64.StdEncoding.DecodeString(cipherText)
	block, err := sm4.NewCipher(key)
	if err != nil {
		panic(err)
	}
	cipherArr := make([]byte, 0)
	cArr := make([]byte, 16)

	j := 0
	for i := 0; i < len(cInArr)/16; i++ {
		original := cInArr[j : j+16]
		block.Decrypt(cArr, original)
		cipherArr = append(cipherArr, cArr...)
		j = j + 16
	}
	originalText := string(PKCS7UnPadding(cipherArr))
	return originalText, nil

}
