package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"github.com/astaxie/beego/logs"
)

const (
	Default_Hex_AESKEY = "1A12B2CBEB9E02D0AF3DB1AD946A9CDF"
	Default_Hex_AESKIV = "B029AC84E67C4A8853FA5E4A50CED3D6"
)

var AES256cbcKey []byte
var AES256cbcIv []byte

func init() {
	AES256cbcKey, _ = hex.DecodeString(Default_Hex_AESKEY)
	AES256cbcIv, _ = hex.DecodeString(Default_Hex_AESKIV)
}

//加密
func EncryptString(s string, aesKey, aesIV string, isPadding bool) string {
	ret := EncryptBytesByKey([]byte(s), aesKey, aesIV, isPadding)
	return string(ret)
}

//加密
func EncryptBytesByKey(s []byte, aesKey, aesIV string, isPadding bool) []byte {

	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		logs.Debug("NewCipher err :%s", err.Error())
		return nil
	}
	cfb := cipher.NewCBCEncrypter(block, []byte(aesIV))
	var plainBytes []byte
	if !isPadding {
		padding := cfb.BlockSize() - len(s)%cfb.BlockSize()
		padding = padding % cfb.BlockSize()
		padtext := bytes.Repeat([]byte{' '}, padding) // make([]byte, padding)
		//fmt.Printf("padding: %d", padding)
		//for i := 0; i < padding; i++ {
		//	padtext[i] = ' '
		//}
		plainBytes = append(s, padtext...)
	} else {

		plainBytes = pKCS7Padding(s, cfb.BlockSize())

	}
	enBytes := make([]byte, len(plainBytes))
	cfb.CryptBlocks(enBytes, plainBytes)
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(enBytes)))
	base64.StdEncoding.Encode(buf, enBytes)
	return buf
}

func pKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	// 填充
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}
