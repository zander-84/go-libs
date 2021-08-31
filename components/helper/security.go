package helper

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

var defaultSecurity = NewSecurity()

type Security struct{}

func NewSecurity() *Security { return new(Security) }

func GetSecurity() *Security { return defaultSecurity }

func (*Security) Sum256(str string) string {
	data := sha256.Sum256([]byte(str))
	return hex.EncodeToString(data[:])
}

func (*Security) Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func (*Security) Sha256Hmac(str string, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(str))
	return hex.EncodeToString(mac.Sum(nil))
}

func (*Security) TripleDESEncrypt(src, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	src = PKCS5Padding(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	dst := src
	blockMode.CryptBlocks(dst, src)
	return dst, nil
}

func (*Security) TripleDESDecrypt(src, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	dst := src
	blockMode.CryptBlocks(dst, src)
	dst = PKCS5UnPadding(dst)
	return dst, nil
}

func (*Security) AESEncrypt(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	src = PKCS5Padding(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	dst := src
	blockMode.CryptBlocks(dst, src)
	return dst, nil
}

func (*Security) AESDecrypt(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	dst := src
	blockMode.CryptBlocks(dst, src)
	dst = PKCS5UnPadding(dst)
	return dst, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - (len(ciphertext) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	newText := append(ciphertext, padText...)
	return newText
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	number := int(origData[length-1])
	return origData[:(length - number)]
}
