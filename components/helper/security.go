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
	"errors"
	"strings"
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

func (s *Security) Md5StrToUpper(str string) string {
	return strings.ToUpper(s.Md5(str))
}

func (*Security) Sha256Hmac(str string, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(str))
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *Security) Sha256HmacToUpper(str string, key string) string {
	return strings.ToUpper(s.Sha256Hmac(str, key))
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

func (*Security) AESCBCPkcs5Encrypt(src, key []byte) ([]byte, error) {
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

func (*Security) AESCBCPkcs5Decrypt(src, key []byte) ([]byte, error) {
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

//AESECBPkcs7Encrypt aes-ecb-encrypt 加密 pkcs7填充
func (*Security) AESECBPkcs7Encrypt(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	src = PKCS7Padding(src, block.BlockSize())
	dst := make([]byte, len(src))
	err = ECBEncrypt(block, dst, src)
	return dst, err
}

//AESECBPkcs7Decrypt aes-ecb-decrypt 解密 pkcs7填充
func (*Security) AESECBPkcs7Decrypt(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, len(src))
	err = ECBDecrypt(block, dst, src)
	if err != nil {
		return nil, err
	}
	return PKCS7UnPadding(dst), nil
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

func PKCS7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS7UnPadding(src []byte) []byte {
	length := len(src)
	if length == 0 {
		return src
	}
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func ECBEncrypt(b cipher.Block, dst, src []byte) error {
	if len(src)%b.BlockSize() != 0 {
		return errors.New("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		return errors.New("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		b.Encrypt(dst, src[:b.BlockSize()])
		src = src[b.BlockSize():]
		dst = dst[b.BlockSize():]
	}
	return nil
}

func ECBDecrypt(b cipher.Block, dst, src []byte) error {
	if len(src)%b.BlockSize() != 0 {
		return errors.New("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		return errors.New("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		b.Decrypt(dst, src[:b.BlockSize()])
		src = src[b.BlockSize():]
		dst = dst[b.BlockSize():]
	}
	return nil
}
