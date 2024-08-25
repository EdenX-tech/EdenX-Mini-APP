package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/spf13/viper"
	"io"
)

var encryptionKey = viper.GetString("secret.encryptionKey")

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid padding size")
	}
	unpadding := int(data[length-1])
	return data[:(length - unpadding)], nil
}

func Encrypt(plainText []byte) (string, error) {
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return "", err
	}

	plainText = pkcs7Padding(plainText, block.BlockSize())

	// 生成初始向量 IV
	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// 使用CBC模式进行加密
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plainText)

	// 返回Base64编码的加密数据
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(cipherText string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cipherText)
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// 去掉 PKCS#7 填充
	return pkcs7Unpadding(ciphertext)

}
