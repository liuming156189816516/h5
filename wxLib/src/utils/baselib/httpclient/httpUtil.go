package httpclient

import (
	"bytes"
	"github.com/astaxie/beego/logs"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/**
获取代理配置
http 代理 //129.226.180.61:8889
socks 代理 129.226.180.61:8889
*/
func GetHttpTransport(sproxy string) *http.Transport {
	if sproxy == "" {
		return nil
	}
	if strings.HasPrefix(sproxy, "//") {
		// http 代理
		proxyUrl, err := url.Parse(sproxy)
		if err != nil {
			logs.Error("Parse http proxy:%s, err:%+v", sproxy, err)
			return nil
		}
		logs.Debug("use http proxy = %s", sproxy)
		return &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	dialer, err := proxy.SOCKS5("tcp", sproxy, nil, proxy.Direct)
	if err != nil {
		logs.Error("proxy.SOCKS5 socks proxy:%s, err:%+v", sproxy, err)
		return nil
	}
	logs.Debug("use socks proxy = %s", sproxy)
	return &http.Transport{Dial: dialer.Dial}
}

// http请求 生成支付的URL
func HttpGet(surl, sproxy string, timeouts ...int64) (error, []byte) {

	request, err := http.NewRequest("GET", surl, nil)
	if err != nil {
		return err, nil
	}
	client := http.Client{}
	if sproxy != "" {
		transPort := GetHttpTransport(sproxy)
		if transPort != nil {
			client.Transport = transPort
		}
		// proxyUrl, err := url.Parse(sproxy)
		// if err != nil {
		// 	logs.LogError("Parse sproxy:%s, err:%+v", sproxy, err)
		// 	return err, nil
		// }
		// logs.LogDebug("sproxy = %s", sproxy)
		// client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	if len(timeouts) > 0 {
		timeout := timeouts[0]
		if timeout > 0 {
			client.Timeout = time.Duration(timeout) * time.Second
		}
	}

	request.Close = true
	resp, err := client.Do(request)
	if err != nil {
		//logs.LogDebug("client.Do %+v err:%+v", request, err)
		logs.Error("HttpProxyRequest method:GET, url:%s, sproxy:%s,  err:%+v, timeout:%+v", surl, sproxy, err, client.Timeout)
		return err, nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Debug("ioutil.ReadAll err:%+v", err)
		return err, nil
	}
	return err, body
}

func HttpProxyRequest(method, surl, sproxy string, body []byte, timeouts ...int64) ([]byte, error) {
	request, err := http.NewRequest(method, surl, bytes.NewReader(body))
	if err != nil {
		logs.Debug("err:%+v", err)
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	if sproxy != "" {
		proxyUrl, err := url.Parse(sproxy)
		if err != nil {
			logs.Error("Parse sproxy:%s, err:%+v", sproxy, err)
			return nil, err
		}
		logs.Debug("sproxy = %s", sproxy)
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	if len(timeouts) > 0 {
		timeout := timeouts[0]
		if timeout > 0 {
			client.Timeout = time.Duration(timeout) * time.Second
		}
	}
	request.Close = true
	resp, err := client.Do(request)
	if err != nil {
		//logs.LogDebug("client.Do %+v err:%+v", request, err)
		logs.Debug("HttpProxyRequest method:%s, url:%s, sproxy:%s, body:%s, err:%+v", method, surl, sproxy, string(body), err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Debug("ioutil.ReadAll err:%+v", err)
		return nil, err
	}
	return body, err
}

func HttpProxyRequestJson(method, surl, sproxy string, body []byte, headerData map[string]string, timeouts ...int64) ([]byte, error) {
	request, err := http.NewRequest(method, surl, strings.NewReader(string(body)))
	if err != nil {
		logs.Debug("err:%+v", err)
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if len(headerData) > 0 {
		for k, v := range headerData {
			request.Header.Set(k, v)
		}
	}
	client := http.Client{}
	if sproxy != "" {
		proxyUrl, err := url.Parse(sproxy)
		if err != nil {
			logs.Error("Parse sproxy:%s, err:%+v", sproxy, err)
			return nil, err
		}
		logs.Debug("sproxy = %s", sproxy)
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	if len(timeouts) > 0 {
		timeout := timeouts[0]
		if timeout > 0 {
			client.Timeout = time.Duration(timeout) * time.Second
		}
	}
	request.Close = true
	resp, err := client.Do(request)
	if err != nil {
		logs.Error("HttpProxyRequest method:%s, sproxy:%s, body:%s, err:%+v", method, surl, sproxy, string(body), err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Debug("ioutil.ReadAll err:%+v", err)
		return nil, err
	}
	return body, err
}

func HttpProxyRequestFormUrlencoded(method, surl, sproxy string, data url.Values, headData map[string]string, timeouts ...int64) ([]byte, error) {
	request, err := http.NewRequest(method, surl, strings.NewReader(data.Encode()))
	if err != nil {
		logs.Debug("err:%+v", err)
		return nil, err
	}

	client := http.Client{}
	if sproxy != "" {
		proxyUrl, err := url.Parse(sproxy)
		if err != nil {
			logs.Error("Parse sproxy sproxy:%s, err:%+v", sproxy, err)
			return nil, err
		}
		logs.Debug("sproxy = %s", sproxy)
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	if len(timeouts) > 0 {
		timeout := timeouts[0]
		client.Timeout = time.Duration(timeout) * time.Second
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	if len(headData) > 0 {
		for k, v := range headData {
			request.Header.Set(k, v)
		}
	}
	request.Close = true
	resp, err := client.Do(request)
	if err != nil {
		logs.Error("HttpProxyRequest method:%s, url:%s, sproxy:%s, data:%+v, err:%+v", method, surl, sproxy, data, err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Debug("ioutil.ReadAll err:%+v", err)
		return nil, err
	}
	return body, err
}
