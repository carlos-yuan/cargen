/**
* @Author: qinwei
* @Date:   2023/3/3 09:07
* @FileUrl:   aes_cbc.go
* @Desc:   AES_CBC 加解密
 */

package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// EncryptCBC5 加密
func EncryptCBC5(plaintext []byte, key, iv string) (string, error) {
	// 创建 AES 加密器
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 填充明文
	padded := PKCS5Padding(plaintext, block.BlockSize())

	// 创建 CBC 加密器

	stream := cipher.NewCBCEncrypter(block, []byte(iv))

	// 加密
	ciphertext := make([]byte, len(padded))
	stream.CryptBlocks(ciphertext, padded)

	// 返回 Base64 编码的密文
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptCBC5 解密
func DecryptCBC5(ciphertext, key, iv string) (str string, err error) {
	defer func() {
		r := recover()
		if r != nil {
			rstr, ok := r.(string)
			if ok {
				err = errors.New(rstr)
			} else {
				rerr, ok := r.(error)
				if ok {
					err = rerr
				}
			}
		}
	}()
	// 解码 Base64 编码的密文
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 创建 AES 解密器
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 创建 CBC 解密器
	stream := cipher.NewCBCDecrypter(block, []byte(iv))

	// 解密
	decrypted := make([]byte, len(ciphertextBytes))
	stream.CryptBlocks(decrypted, ciphertextBytes)

	// 去除填充
	unPadded := PKCS5Trimming(decrypted)
	// 返回明文
	return string(unPadded), nil
}

// EncryptCBC5Bytes 加密
func EncryptCBC5Bytes(plaintext []byte, key, iv string) ([]byte, error) {
	// 创建 AES 加密器
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	// 填充明文
	padded := PKCS5Padding(plaintext, block.BlockSize())

	// 创建 CBC 加密器

	stream := cipher.NewCBCEncrypter(block, []byte(iv))

	// 加密
	ciphertext := make([]byte, len(padded))
	stream.CryptBlocks(ciphertext, padded)

	// 返回 Base64 编码的密文
	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// DecryptCBC5Bytes 解密
func DecryptCBC5Bytes(ciphertext, key, iv string) (str []byte, err error) {
	defer func() {
		r := recover()
		if r != nil {
			rstr, ok := r.(string)
			if ok {
				err = errors.New(rstr)
			} else {
				rerr, ok := r.(error)
				if ok {
					err = rerr
				}
			}
		}
	}()
	// 解码 Base64 编码的密文
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	// 创建 AES 解密器
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	// 创建 CBC 解密器
	stream := cipher.NewCBCDecrypter(block, []byte(iv))

	// 解密
	decrypted := make([]byte, len(ciphertextBytes))
	stream.CryptBlocks(decrypted, ciphertextBytes)

	// 去除填充
	unPadded := PKCS5Trimming(decrypted)
	// 返回明文
	return unPadded, nil
}

// EncryptCBC7 加密
func EncryptCBC7(plaintext []byte, key []byte) (string, error) {
	// 创建 AES 加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 填充明文
	padded := pkcs7Padding(plaintext, block.BlockSize())

	// 创建 CBC 加密器
	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCBCEncrypter(block, iv)

	// 加密
	ciphertext := make([]byte, len(padded))
	stream.CryptBlocks(ciphertext, padded)

	// 返回 Base64 编码的密文
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptCBC7 解密
func DecryptCBC7(ciphertext string, key []byte) (string, error) {
	// 解码 Base64 编码的密文
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 创建 AES 解密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 创建 CBC 解密器
	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCBCDecrypter(block, iv)

	// 解密
	decrypted := make([]byte, len(ciphertextBytes))
	stream.CryptBlocks(decrypted, ciphertextBytes)

	// 去除填充
	unPadded, err := pkcs7UnPadding(decrypted)
	if err != nil {
		return "", err
	}

	// 返回明文
	return string(unPadded), nil
}
