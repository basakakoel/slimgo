package context

import (
	"net/http"
	"strconv"
	"strings"
)

type ContextInput struct {
	Request *http.Request
}

func NewContextInput(req *http.Request) *ContextInput {
	return &ContextInput{
		Request: req,
	}
}

//Protocol
func (this *ContextInput) Protocol() string {
	return this.Request.Proto
}

//Url
func (this *ContextInput) Url() string {
	return this.Request.URL.Path
}

//Uri
func (this *ContextInput) Uri() string {
	return this.Request.RequestURI
}

//Host
func (this *ContextInput) Host() string {
	if this.Request.Host != "" {
		hostArr := strings.Split(this.Request.Host, ":")
		if len(hostArr) > 0 {
			return hostArr[0]
		}
		return this.Request.Host
	}
	return "localhost"
}

//Method
func (this *ContextInput) IsMethod(method string) bool {
	return this.Request.Method == method
}

func (this *ContextInput) IsGet() bool {
	return this.IsMethod("GET")
}

func (this *ContextInput) IsPost() bool {
	return this.IsMethod("POST")
}

func (this *ContextInput) IsPut() bool {
	return this.IsMethod("PUT")
}

func (this *ContextInput) IsDelete() bool {
	return this.IsMethod("DELETE")
}

func (this *ContextInput) IsPatch() bool {
	return this.IsMethod("PATCH")
}

//Header
func (this *ContextInput) Header(key string) string {
	return this.Request.Header.Get(key)
}

//Is upload form
func (this *ContextInput) IsUpload() bool {
	return strings.Contains(this.Header("Content-Type"), "multipart/form-data")
}

//Is websocket
func (this *ContextInput) IsWebsocket() bool {
	return this.Header("Upgrade") == "websocket"
}

//Proxy
func (this *ContextInput) Proxy() []string {
	if ips := this.Header("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

//IP
func (this *ContextInput) Ip() string {
	ips := this.Proxy()
	if len(ips) > 0 && ips[0] != "" {
		ipsArr := strings.Split(ips[0], ":")
		return ipsArr[0]
	}
	ipArr := strings.Split(this.Request.RemoteAddr, ":")
	if len(ipArr) > 0 {
		if ipArr[0] != "[" {
			return ipArr[0]
		}
	}
	return "127.0.0.1"
}

//Port
func (this *ContextInput) Port() int {
	portArr := strings.Split(this.Request.Host, ":")
	if len(portArr) > 1 {
		port, _ := strconv.Atoi(portArr[1])
		return port
	}
	return 80
}

//Referer
func (this *ContextInput) Referer() string {
	return this.Header("Referer")
}

//Cookie
func (this *ContextInput) Cookie(key string) string {
	cookie, err := this.Request.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}

//get param
func (this *ContextInput) GetParam(key string) string {
	if this.Request.Form == nil {
		this.Request.ParseForm()
	}
	return this.Request.Form.Get(key)
}
