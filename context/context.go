package context

import (
	"net/http"
)

type Context struct {
	Input          *ContextInput
	Output         *ContextOutput
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

//redirect
func (this *Context) Redirect(url string, status int) {
	this.Output.Header("Location", url)
	this.ResponseWriter.WriteHeader(status)
}

//write
func (this *Context) WriteString(content string) {
	this.ResponseWriter.Write([]byte(content))
}

func (this *Context) GetCookie(key string) string {
	return this.Input.Cookie(key)
}

func (this *Context) SetCookie(key, value string, params ...interface{}) {
	this.Output.Cookie(key, value, params...)
}
