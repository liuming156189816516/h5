package wxHttp

import (
	"crypto/tls"
	"errors"
	"github.com/go-resty/resty"
	"net/http"
	"strings"
	"time"
)

type ZHttpReqParam struct {
	Url        string
	Method     string
	Content    interface{}
	Headers    map[string]string
	FormData   map[string]string
	Proxy      string
	Timeout    time.Duration
	Transport  *http.Transport
	NORedirect bool
	Debug      bool
	Trace      bool
}

type ZHttpRespParam struct {
	Err       error
	Status    int
	Cookies   []*http.Cookie
	Body      []byte
	Headers   http.Header
	TraceInfo ZHttpRespTraceInfo
}

type ZHttpRespTraceInfo struct {
	DNSLookup      time.Duration
	ConnTime       time.Duration
	TCPConnTime    time.Duration
	TLSHandshake   time.Duration
	ServerTime     time.Duration
	ResponseTime   time.Duration
	TotalTime      time.Duration
	IsConnReused   bool
	IsConnWasIdle  bool
	ConnIdleTime   time.Duration
	RequestAttempt int
	RemoteAddr     string
}

func ZHttp(Req ZHttpReqParam) ZHttpRespParam {

	Http := resty.New()

	if Req.Transport != nil {
		Http.SetTransport(Req.Transport)
	} else {
		Http.SetTransport(&http.Transport{
			MaxIdleConnsPerHost: -1,
			DisableKeepAlives:   true,
			ForceAttemptHTTP2:   false,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		})
	}

	if Req.Timeout == time.Duration(0) {
		Req.Timeout = time.Duration(30)
	}

	Http.SetTimeout(Req.Timeout * time.Second)

	Http.SetDebug(Req.Debug)

	if Req.Headers != nil {
		Http.SetHeaders(Req.Headers)
	}

	if Req.Proxy != "" {
		Http.SetProxy(Req.Proxy)
	}

	if Req.Trace {
		Http.EnableTrace()
	}

	//禁止重定向
	if Req.NORedirect {
		Http.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}))
	}

	var request *resty.Response
	var err error

	switch strings.ToLower(Req.Method) {
	case strings.ToLower("POST"):
		request, err = Http.R().SetBody(Req.Content).Post(Req.Url)
	case strings.ToLower("GET"):
		request, err = Http.R().SetBody(Req.Content).Get(Req.Url)
	case strings.ToLower("postForm"):
		request, err = Http.R().SetFormData(Req.FormData).Post(Req.Url)
	default:
		return ZHttpRespParam{Err: errors.New("不存在的请求类型。")}
	}
	if err != nil {
		return ZHttpRespParam{Err: err}
	}

	TraceInfo := ZHttpRespTraceInfo{}

	if request.Request != nil {
		//time
		ti := request.Request.TraceInfo()
		TraceInfo = ZHttpRespTraceInfo{
			DNSLookup:      ti.DNSLookup,
			ConnTime:       ti.ConnTime,
			TCPConnTime:    ti.TCPConnTime,
			TLSHandshake:   ti.TLSHandshake,
			ServerTime:     ti.ServerTime,
			ResponseTime:   ti.ResponseTime,
			TotalTime:      ti.TotalTime,
			IsConnReused:   ti.IsConnReused,
			IsConnWasIdle:  ti.IsConnWasIdle,
			ConnIdleTime:   ti.ConnIdleTime,
			RequestAttempt: ti.RequestAttempt,
		}

		if ti.RemoteAddr != nil {
			TraceInfo.RemoteAddr = ti.RemoteAddr.String()
		}
	}

	return ZHttpRespParam{
		Err:       err,
		Body:      request.Body(),
		Headers:   request.Header(),
		Status:    request.StatusCode(),
		Cookies:   request.Cookies(),
		TraceInfo: TraceInfo,
	}
}
