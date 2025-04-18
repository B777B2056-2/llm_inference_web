package proxy

import (
	"fmt"
	"llm_online_interence/llmgateway/confparser"
	"llm_online_interence/llmgateway/resource"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorHandlerType = func(w http.ResponseWriter, r *http.Request, err error)

// proxyErrorHandler 反向代理错误处理函数
func proxyErrorHandler(proxy *httputil.ReverseProxy, svcName, proxyPath string) errorHandlerType {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		if err == nil {
			return
		}

		requestDump, _ := httputil.DumpRequest(r, true)
		resource.Logger.WithFields(logrus.Fields{
			"backend": svcName,
			"path":    proxyPath,
			"request": string(requestDump),
			"error":   err.Error(),
		}).Error("proxy request error")
		w.WriteHeader(http.StatusBadGateway)
	}
}

// commonProxyHandler 通用反向代理处理函数
func commonProxyHandler(ctx *gin.Context, backend confparser.BackendConfigItem) {
	protocol := backend.Protocol
	svcName := backend.SvcName
	connTimeout := backend.ConnectionTimeout
	respTimeout := backend.ResponseTimeout

	// 初始化反向代理器
	remote, err := url.Parse(fmt.Sprintf("%s://%s", protocol, svcName))
	if err != nil {
		panic(err)
	}

	proxyPath := ctx.Param("proxyPath")
	proxy := httputil.NewSingleHostReverseProxy(remote)

	requestDump, _ := httputil.DumpRequest(ctx.Request, true)
	resource.Logger.WithFields(logrus.Fields{
		"backend": backend.SvcName,
		"path":    proxyPath,
		"request": string(requestDump),
	}).Debug("proxy request")

	// 设置反向代理的请求转发规则
	proxy.Director = func(req *http.Request) {
		req.Header = ctx.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = proxyPath
	}

	// 设置连接和响应超时
	proxy.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Duration(connTimeout) * time.Millisecond,
		}).DialContext,
		ResponseHeaderTimeout: time.Duration(respTimeout) * time.Millisecond,
	}

	// 绑定错误处理函数
	proxy.ErrorHandler = proxyErrorHandler(proxy, svcName, proxyPath)

	// 执行反向代理
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
