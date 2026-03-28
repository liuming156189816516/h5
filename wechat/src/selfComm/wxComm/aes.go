package wxComm

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	jsoniter "github.com/json-iterator/go"
)

// 明文补码算法
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 明文减码算法
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AES CBC 加密
func AesEncryptCbc(rsp interface{}) string {
	ret := ""
	key := "8dw/JfjjoMs0dzVGOX2ntb1iw2k9+JD4"
	iv := "ZGdIobme/Sb4Idwg"
	keyByte := []byte(key)
	ivByte := []byte(iv)
	origData, err := jsoniter.Marshal(rsp)
	if err != nil {
		return ret
	}
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return ret
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)

	//获取CBC加密模式
	mode := cipher.NewCBCEncrypter(block, ivByte)
	crypted := make([]byte, len(origData))
	mode.CryptBlocks(crypted, origData)
	ret = base64.StdEncoding.EncodeToString(crypted)
	return ret
}

// AES CBC 解密
func AesDecryptCbc(requestBody []byte) []byte {
	key := "8dw/JfjjoMs0dzVGOX2ntb1iw2k9+JD4"
	iv := "ZGdIobme/Sb4Idwg"
	ret := []byte{}
	body := &RequestBody{}
	err := jsoniter.Unmarshal(requestBody, &body)
	if err != nil {
		return ret
	}
	encryption, err := base64.StdEncoding.DecodeString(body.Key)
	if err != nil {
		return ret
	}
	crypted := encryption
	keyByte := []byte(key)
	ivByte := []byte(iv)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return ret
	}
	mode := cipher.NewCBCDecrypter(block, ivByte)
	plaintext := make([]byte, len(crypted))
	mode.CryptBlocks(plaintext, crypted)
	ret = PKCS5UnPadding(plaintext)
	return ret
}

type RequestBody struct {
	Key string `json:"key"`
}
