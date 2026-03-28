package comm

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 单独的request请求 - json提交
func HttpProxyRequest(method, surl, sproxy string, body []byte, timeouts ...int64) ([]byte, error) {
	request, err := http.NewRequest(method, surl, bytes.NewReader(body))
	if err != nil {
		logs.Error("url:%s, proxy:%s, err:%+v", surl, sproxy, err)
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	if sproxy != "" {
		proxyUrl, err := url.Parse(sproxy)
		if err != nil {
			logs.Trace("Parse sproxy err:%+v", err)
			return nil, err
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	if len(timeouts) > 0 {
		timeout := timeouts[0]
		client.Timeout = time.Duration(timeout) * time.Second
	}
	request.Close = true
	resp, err := client.Do(request)
	if err != nil {
		logs.Error("client.Do url:%s, req:%+v, err:%+v", surl, request, err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Trace("ioutil.ReadAll err:%+v", err)
		return nil, err
	}
	return body, err
}

// 单独的request请求 - 表单提交
func HttpProxyRequestByForm(method, surl, sproxy string, body []byte, timeouts ...int64) ([]byte, error) {
	var tempMap map[string]string
	err := json.Unmarshal(body, &tempMap)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for k, v := range tempMap {
		writer.WriteField(k, v)
	}
	err = writer.Close()
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	request, err := http.NewRequest(method, surl, bytes.NewReader(body))
	if err != nil {
		logs.Error("url:%s, proxy:%s, err:%+v", surl, sproxy, err)
		return nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	client := http.Client{}
	if sproxy != "" {
		proxyUrl, err := url.Parse(sproxy)
		if err != nil {
			logs.Trace("Parse sproxy err:%+v", err)
			return nil, err
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	if len(timeouts) > 0 {
		timeout := timeouts[0]
		client.Timeout = time.Duration(timeout) * time.Second
	}
	request.Close = true
	resp, err := client.Do(request)
	if err != nil {
		logs.Error("client.Do url:%s, req:%+v, err:%+v", surl, request, err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Trace("ioutil.ReadAll err:%+v", err)
		return nil, err
	}
	return body, err
}

// 垃圾四方专用 - x-www-form-urlencoded
func HttpProxyRequestBySF(method, surl, sproxy string, data string, timeouts ...int64) bool {
	request, err := http.NewRequest(method, surl, strings.NewReader(data))
	if err != nil {
		return false
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	logs.Info("请求四方返回值resp====>", resp)
	logs.Info("请求四方返回值err====>", err)
	if err != nil {
		return false
	}
	if len(timeouts) > 0 {
		timeout := timeouts[0]
		client.Timeout = time.Duration(timeout) * time.Second
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}
}

func HttpProxyRequestBySFGetCode(url string, timeouts ...int64) string {
	//url := "https://sms.xn--ssl-8e5fl55p.com/api/smsforjson?Room=v6msOOqGKI%0A"
	method := "POST"
	client := &http.Client{
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return ""
	}
	res, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	return string(body)

}
