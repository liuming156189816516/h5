package httpclient

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	GlobalHttpClientTimeOut = 3 //GlobalHttpClientTimeOut is a general config

	constContentTypeJson           = 1
	constContentTypeFormUrlencoded = 2
)


func DoHttpClientPost(clt *http.Client, reqUrl string, postBody string, timeout int,cookies string, contentType int,headers ...string) (body string, err error){
	if timeout == 0 {
		timeout = GlobalHttpClientTimeOut
	}
	//NewRequest
	reqest, errReq := http.NewRequest("POST", reqUrl, strings.NewReader(postBody))
	if errReq != nil {
		return "", errReq
	}
	//cookie
	reqest.Header.Set("Cookie", cookies)
	for i := 0; i < len(headers)-1; {
		reqest.Header.Set(headers[i], headers[i+1])
		i += 2
	}
	if constContentTypeFormUrlencoded == contentType {
		reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		reqest.Header.Set("Content-Type", "application/json")
	}
	reqest.Header.Set("Accept", "*/*")

	//do send()
	resp, errDo := clt.Do(reqest)
	if errDo != nil {
		err = errDo
	} else {
		result, errRead := ioutil.ReadAll(resp.Body)
		if errRead != nil {
			err = errRead
		}
		resp.Body.Close()
		body = string(result)
	}

	return body, err
}

func innerDoHttpPostWithTimeout(uin uint64, reqUrl string,  postBody string, timeout int,cookies string, contentType int, headers ...string) (body string, err error) {
	//new client
	if timeout == 0 {
		timeout = GlobalHttpClientTimeOut
	}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, errDial := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout)) //connect timeout
				if errDial != nil {
					return nil, errDial
				}
				errDeadLine := conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout)))
				if errDeadLine != nil {
					return nil, errDeadLine
				}
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * time.Duration(timeout), //sets the read timeout
		},
	}
	return DoHttpClientPost(client, reqUrl, postBody,timeout, cookies,  contentType, headers...)

}

// DoHttpPostForm
func DoHttpPostFormHeader(uin uint64, reqUrl string, postBody string, cookies string) (body string, err error) {

	return innerDoHttpPostWithTimeout(uin, reqUrl, postBody, 0, cookies, constContentTypeFormUrlencoded)
}

// DoHttpPost is a wrapper
func DoHttpPost(uin uint64, reqUrl string, postBody string, cookies string, headers ...string) (body string, err error) {

	return innerDoHttpPostWithTimeout(uin, reqUrl, postBody,0,  cookies, constContentTypeJson, headers...)
}
func DoHttpPostWithTimeout(uin uint64, reqUrl string, postBody string, timeout int, cookies string, headers ...string) (body string, err error) {

	return innerDoHttpPostWithTimeout(uin, reqUrl,  postBody, timeout,cookies, constContentTypeJson, headers...)
}

// DoHttpGet is a wrapper
func DoHttpGet(uin uint64, reqUrl string, cookies string) (body string, err error) {
	return DoHttpGetEx(uin, reqUrl, cookies, GlobalHttpClientTimeOut)
}

func DoHttpGetEx(uin uint64, reqUrl string, cookies string, timeout int, headers ...string) (body string, err error) {
	//new client
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, errDial := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout)) //connect timeout
				if errDial != nil {
					return nil, errDial
				}
				errDeadLine := conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout)))
				if errDeadLine != nil {
					return nil, errDeadLine
				}
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * time.Duration(timeout), //sets the read timeout
		},
	}

	reqest, errReq := http.NewRequest("GET", reqUrl, nil)
	if errReq != nil {
		return "", errReq
	}
	reqest.Header.Set("Cookie", cookies)
	for i := 0; i < len(headers)-1; {
		reqest.Header.Set(headers[i], headers[i+1])
		i += 2
	}
	resp, errDo := client.Do(reqest)
	if errDo != nil {
		err = errDo
	} else {
		result, errRead := ioutil.ReadAll(resp.Body)
		if errRead != nil {
			err = errRead
		}
		resp.Body.Close()
		body = string(result)
	}

	return body, err
}
