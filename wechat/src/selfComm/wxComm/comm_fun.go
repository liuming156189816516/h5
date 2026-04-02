package wxComm

import (
	"bytes"
	"comm/comm"
	"comm/cos"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	jsoniter "github.com/json-iterator/go"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"selfComm/wxComm/wxHttp"
	"utils"
)

type GetRealIpRsp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Addr        string `json:"addr"`
		Country     string `json:"country"`
		Isp         string `json:"isp"`
		CountryCode string `json:"country_code"`
	} `json:"data"`
}

// 获取真实ip
func GetRealIp(proxyType, username, password, address string) bool {
	proxyIp := fmt.Sprintf(proxyType+"://%s:%s@%s", username, password, address)
	ret := &GetRealIpRsp{}
	ret.Code = -1
	rsp := wxHttp.ZHttp(wxHttp.ZHttpReqParam{
		Url:     "https://api.ip.sb/geoip?callback=getgeoip",
		Method:  "get",
		Timeout: 10,
		Proxy:   proxyIp,
	})
	if rsp.Status == 200 {
		return true

	}
	return false
}

type GetIpInfoRsp struct {
	Status  string `json:"status"`
	Country string `json:"country"`
}

// 获取ip信息
func GetIpInfo(ip string) *GetIpInfoRsp {
	apiUrl := "http://ip-api.com/json/" + ip + "?lang=zh-EN"
	ret := &GetIpInfoRsp{}
	rsp := wxHttp.ZHttp(wxHttp.ZHttpReqParam{
		Url:     apiUrl,
		Method:  "get",
		Timeout: 3,
	})
	if rsp.Status == 200 {
		jsoniter.Unmarshal(rsp.Body, ret)
		return ret
	}
	return ret
}

func GetAccountFriendId(sendAccount, friendAccount string) string {
	accFriendId := ""
	send := utils.StrToInt64(sendAccount)
	friend := utils.StrToInt64(friendAccount)
	if send > friend {
		accFriendId = sendAccount + "_" + friendAccount
	} else {
		accFriendId = friendAccount + "_" + sendAccount
	}
	return accFriendId
}

func GetBase64ByUrl(url string) (base64Str, w, h string) {
	resp, err := http.Get(url)
	if err != nil {
		return "", "", ""
	}
	defer resp.Body.Close()
	all, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return "", "", ""
	}
	buf := bytes.NewBuffer(all)
	config, _, _ := image.DecodeConfig(buf)
	w1 := utils.IntToStr(int64(config.Width))
	h1 := utils.IntToStr(int64(config.Height))
	urlbase64 := base64.StdEncoding.EncodeToString(all)
	return urlbase64, w1, h1
}

func GetBase64ByUrlSmall(url string, w1, h1 int) string {
	str := ""
	resp, err := http.Get(url)
	if err != nil {
		return str
	}
	defer resp.Body.Close()
	imgData, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return str
	}
	buf := bytes.NewBuffer(imgData)
	m, _, err := image.Decode(buf) // 图片文件解码
	if m == nil {
		return str
	}
	nw1 := 50
	nh1 := 50
	for i := 2; i < 100; i++ {
		nw := w1 / i
		nh := h1 / i
		if nw <= 100 && nh <= 100 {
			nw1 = nw
			nh1 = nh
			break
		}
	}
	subImg := resize.Resize(uint(nw1), uint(nh1), m, resize.Lanczos3)
	emptyBuff := bytes.NewBuffer(nil) //开辟一个新的空buff
	err = jpeg.Encode(emptyBuff, subImg, nil)
	if err == nil {
		str = base64.StdEncoding.EncodeToString(emptyBuff.Bytes())
	}
	return str
}

func GetBase64ByUrlSmall2(url string, w1, h1 int) string {
	str := ""
	resp, err := http.Get(url)
	if err != nil {
		return str
	}
	defer resp.Body.Close()
	imgData, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return str
	}
	buf := bytes.NewBuffer(imgData)
	m, _, err := image.Decode(buf) // 图片文件解码
	if m == nil {
		return str
	}
	subImg := resize.Resize(uint(200), uint(200), m, resize.Lanczos3)
	emptyBuff := bytes.NewBuffer(nil) //开辟一个新的空buff
	err = jpeg.Encode(emptyBuff, subImg, nil)
	if err == nil {
		str = base64.StdEncoding.EncodeToString(emptyBuff.Bytes())
	}
	return str
}

type TranslateTextRsp struct {
	Data struct {
		Translations []Translation `json:"translations"`
	} `json:"data"`
}

type Translation struct {
	TranslatedText         string `json:"translatedText"`
	DetectedSourceLanguage string `json:"detectedSourceLanguage"`
}

func TranslateText(targetLanguage, text string) string {
	text = url.QueryEscape(text)
	rawUrl := "https://translation.googleapis.com/language/translate/v2?key=AIzaSyADyxGF3Fvhm6XE6h0pJauxI3wCO5srBAk" + "&q=" + text + "&target=" + targetLanguage + "&format=text"
	rsp := wxHttp.ZHttp(wxHttp.ZHttpReqParam{
		Url:     rawUrl,
		Method:  "get",
		Timeout: 10,
	})
	ret := &TranslateTextRsp{}
	jsoniter.UnmarshalFromString(string(rsp.Body), &ret)
	if len(ret.Data.Translations) > 0 {
		return ret.Data.Translations[0].TranslatedText
	}
	return ""
}

type TinyUrlPara struct {
	Url    string `json:"url"`
	Domain string `json:"domain"`
	//Alias       string `json:"alias"`
	//Tags        string `json:"tags"`
	//ExpiresAt   string `json:"expires_at"`
	//Description string `json:"description"`
}
type TinyUrlRsp struct {
	Code int64 `json:"code"`
	Data struct {
		Domain  string `json:"domain"`
		TinyUrl string `json:"tiny_url"`
	} `json:"data"`
}

func TinyUrl(url string) string {
	api := "https://api.tinyurl.com/create?api_token=g3epVxM1jcDFXfCljvM5dNmPP5L2dGAA4p3GZDLoo3y2mC4bSEwpO7hHk1cF"
	req := &TinyUrlPara{}
	req.Url = url
	req.Domain = "tinyurl.com"
	rsp := wxHttp.ZHttp(wxHttp.ZHttpReqParam{
		Url:     api,
		Method:  "post",
		Content: req,
		Timeout: 10,
	})
	ret := &TinyUrlRsp{}
	jsoniter.UnmarshalFromString(string(rsp.Body), &ret)
	if ret.Code == 0 {
		return ret.Data.TinyUrl
	}
	return ""
}

func IsGd(url2 string) string {
	encodedURL := url.QueryEscape(url2)
	api := "https://is.gd/create.php?format=simple&url=" + encodedURL
	rsp := wxHttp.ZHttp(wxHttp.ZHttpReqParam{
		Url:     api,
		Method:  "get",
		Timeout: 10,
	})
	if rsp.Status == 200 {
		return string(rsp.Body)
	}
	return ""
}

func ToFile(content, ctype string) string {
	retStr := ""
	tmpPath := beego.AppConfig.String("tmpPath")
	fileName := ""
	if ctype == "image" {
		fileName = comm.Md5(content) + ".jpg"
	} else if ctype == "ptt" || ctype == "audio" {
		fileName = comm.Md5(content) + ".mp3"
	} else if ctype == "video" {
		fileName = comm.Md5(content) + ".mp4"
	}
	filePath := tmpPath + fileName

	ddd, _ := base64.StdEncoding.DecodeString(content) //成图片文件并把文件写入到buffer
	err2 := ioutil.WriteFile(filePath, ddd, 0777)      //buffer输出到jpg文件中（不做处理，直接写到文件）
	if err2 != nil {
		return retStr
	}
	defer os.Remove(filePath)
	fileUrl := cos.UploadAwsFile(filePath, fileName)
	return fileUrl
}

func IsChain(phone string) bool {
	//if strings.HasPrefix(phone, "+852") || strings.HasPrefix(phone, "852") || strings.HasPrefix(phone, "00852") ||
	//	strings.HasPrefix(phone, "+86") || strings.HasPrefix(phone, "86") || strings.HasPrefix(phone, "0086") ||
	//	strings.HasPrefix(phone, "+886") || strings.HasPrefix(phone, "886") || strings.HasPrefix(phone, "00886") ||
	//	strings.HasPrefix(phone, "+853") || strings.HasPrefix(phone, "853") || strings.HasPrefix(phone, "00853") {
	//	return true
	//}
	return false
}

func JfifEncode(w io.Writer, m image.Image, o *jpeg.Options) error {
	return jpeg.Encode(&jfifWriter{w: w}, m, o)
}

// jfifWriter wraps an io.Writer to convert the data written to it from a plain
// JPEG to a JFIF-enhanced JPEG. It implicitly buffers the first three bytes
// written to it. The fourth byte will tell whether the original JPEG already
// has the APP0 chunk that JFIF requires.
type jfifWriter struct {
	// w is the wrapped io.Writer.
	w io.Writer
	// n ranges between 0 and 4 inclusive. It is the number of bytes written to
	// this (which also implements io.Writer), saturating at 4. The first three
	// bytes are expected to be {0xff, 0xd8, 0xff}. The fourth byte indicates
	// whether the second JPEG chunk is an APP0 chunk or something else.
	n int
}

func (jw *jfifWriter) Write(p []byte) (int, error) {
	nSkipped := 0

	for jw.n < 3 {
		if len(p) == 0 {
			return nSkipped, nil
		} else if p[0] != jfifChunk[jw.n] {
			return nSkipped, errors.New("jfifWriter: input was not a JPEG")
		}
		nSkipped++
		jw.n++
		p = p[1:]
	}

	if jw.n == 3 {
		if len(p) == 0 {
			return nSkipped, nil
		}
		chunk := jfifChunk
		if p[0] == 0xe0 {
			// The input JPEG already has an APP0 marker. Just write SOI (2
			// bytes) and an 0xff: the three bytes we've previously skipped.
			chunk = chunk[:3]
		}
		if _, err := jw.w.Write(chunk); err != nil {
			return nSkipped, err
		}
		jw.n = 4
	}

	n, err := jw.w.Write(p)
	return n + nSkipped, err
}

// jfifChunk is a sequence: an SOI chunk, an APP0/JFIF chunk and finally the
// 0xff that starts the third chunk.
var jfifChunk = []byte{
	0xff, 0xd8, // SOI  marker.
	0xff, 0xe0, // APP0 marker.
	0x00, 0x10, // Length: 16 byte payload (including these two bytes).
	0x4a, 0x46, 0x49, 0x46, 0x00, // "JFIF\x00".
	0x01, 0x01, // Version 1.01.
	0x00,       // No density units.
	0x00, 0x01, // Horizontal pixel density.
	0x00, 0x01, // Vertical   pixel density.
	0x00, // Thumbnail width.
	0x00, // Thumbnail height.
	0xff, // Start of the third chunk's marker.
}

type Advertise struct {
	Title   string `json:"title"`   //标题
	Img     string `json:"img"`     //图片
	Url     string `json:"url"`     //链接
	Remark  string `json:"remark"`  //描述
	Content string `json:"content"` //内容
	IsShow  bool   `json:"isShow"`  //展示广告
}
