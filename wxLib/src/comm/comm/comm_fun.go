package comm

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/json-iterator/go"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	zip2 "github.com/mzky/zip"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const PwdSalt = "d7a7447682443798dskjhsfd9238fa0932baec9"

//手机号码 key
func GetPhoneKey(code int64, phone string) string {
	return fmt.Sprintf("+%d %s", code, phone)
}

// 密码加盐
func GetSaltPwd(pwd string) string {
	h := md5.New()
	h.Write([]byte(pwd + PwdSalt))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//md5
func Md5ByInterface(data interface{}) string {
	str, ok := data.(string)
	if !ok {
		str, _ = jsoniter.MarshalToString(data)
	}
	h := md5.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//md5
func Md5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//md5
func MD5ToLower(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//二进制转换md5
func Md5FromBin(data []byte) string {
	h := md5.New()
	h.Write(data)
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

var baseStr string = "0123456789ABCDEFGHJKLMNPQRSTUVWXYZ"

func GetCode() string {
	ll := len(baseStr)
	res := []byte{}
	rand.Seed(time.Now().UnixNano())
	for i := int32(0); i < 8; i++ {
		r := rand.Int() % ll
		res = append(res, baseStr[r])
	}
	return string(res)
}

func CheckInStr(in string, list []string) bool {
	for _, c := range list {
		if c == in {
			return true
		}
	}
	return false
}
func UseNewEncoderGbk(str string) string {
	enc := mahonia.NewEncoder("gbk")
	ret := enc.ConvertString(str)
	//if !ok {
	//	logs.Error("解析错误 %+v", str)
	//	return str
	//}
	return ret
}

func ZipFiles(filename string, files []string, oldform, newform string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// 把files添加到zip中
	for _, file := range files {

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// 获取file的基础信息
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		//使用上面的FileInforHeader() 就可以把文件保存的路径替换成我们自己想要的了，如下面
		header.Name = strings.Replace(file, oldform, newform, -1)

		// 优化压缩
		// 更多参考see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, zipfile); err != nil {
			return err
		}
	}
	return nil
}

func IsZip(zipPath string) bool {
	f, err := os.Open(zipPath)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 4)
	if n, err := f.Read(buf); err != nil || n < 4 {
		return false
	}

	return bytes.Equal(buf, []byte("PK\x03\x04"))
}

// password值可以为空""
func ZipPw(zipPath, password string, fileList []string, oldform, newform string) error {
	if len(fileList) < 1 {
		return fmt.Errorf("将要压缩的文件列表不能为空")
	}
	fz, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	zw := zip2.NewWriter(fz)
	defer zw.Close()

	for _, fileName := range fileList {
		fr, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer fr.Close()

		// 获取file的基础信息
		info, err := fr.Stat()
		if err != nil {
			return err
		}

		header, err := zip2.FileInfoHeader(info)
		if err != nil {
			return err
		}
		fileName = strings.Replace(fileName, oldform, newform, -1)
		//防止中文乱码
		//fileName = url.QueryEscape(fileName)
		//使用上面的FileInforHeader() 就可以把文件保存的路径替换成我们自己想要的了，如下面
		header.Name = fileName

		// 优化压缩
		// 更多参考see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip2.Deflate
		header.SetPassword(password)
		header.SetEncryptionMethod(zip2.AES256Encryption)

		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, fr); err != nil {
			return err
		}
	}
	return zw.Flush()
}

// password值可以为空""
// 当decompressPath值为"./"时，解压到相对路径
func UnZipPw(zipPath, password, decompressPath string) error {
	if !IsZip(zipPath) {
		return fmt.Errorf("压缩文件格式不正确或已损坏")
	}
	r, err := zip2.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if password != "" {
			if f.IsEncrypted() {
				f.SetPassword(password)
			} else {
				return errors.New("must be encrypted")
			}
		}
		fp := filepath.Join(decompressPath, f.Name)
		dir, _ := filepath.Split(fp)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}

		w, err := os.Create(fp)
		if nil != err {
			return err
		}

		fr, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(w, fr)
		if err != nil {
			return err
		}
		w.Close()
	}
	return nil
}

//生成指定位数随机字符串
func GetRandomString(lenght int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	bytesLen := len(bytes)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lenght; i++ {
		result = append(result, bytes[r.Intn(bytesLen)])
	}
	return string(result)
}

//  解析二维码--base64转url
func Base64ToUrl(base string) string {
	url := ""
	data64, err := base64.StdEncoding.DecodeString(base)
	if err != nil {
		return url
	}
	fi := strings.NewReader(string(data64))
	img, _, err := image.Decode(fi)
	if err == nil {
		bmp, err := gozxing.NewBinaryBitmapFromImage(img)
		if err == nil {
			qrReader := qrcode.NewQRCodeReader()
			result, err := qrReader.Decode(bmp, nil)
			if err == nil {
				url = result.String()
			}
		}
	}
	return url
}

func Get62Data(str string) string {
	DevicelId := MD5Encode(strings.ToUpper(MD5Encode(str)))
	return "62706c6973743030d4010203040506090a582476657273696f6e58246f626a65637473592461726368697665725424746f7012000186a0a2070855246e756c6c5f1020" + hex.EncodeToString([]byte(DevicelId)) + "5f100f4e534b657965644172636869766572d10b0c54726f6f74800108111a232d32373a406375787d0000000000000101000000000000000d0000000000000000000000000000007f"
}

func MD5Encode(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//巴基斯坦 +92
//俄罗斯 +7
//马来西亚  +60
//美国  +1
//摩洛哥 +212
//印度尼西亚 +62
//越南 +84
//中国    +86
//香港 +852
var country = map[string]string{"+92": "PK", "+7": "RU", "+60": "MY", "+1": "US", "+212": "MA", "+62": "ID", "+84": "VN", "+852": "HK"}

//国家代码
func GetCountryAbbreviation(str string) string {
	if _, ok := country[str]; ok {
		return country[str]
	}
	return "RU"
}
