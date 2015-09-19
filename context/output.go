package context

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

type ContextOutput struct {
	Context *Context
	Status  int
}

func NewContextOutput() *ContextOutput {
	return &ContextOutput{}
}

//Header
func (this *ContextOutput) Header(key, value string) {
	this.Context.ResponseWriter.Header().Set(key, value)
}

//Body
func (this *ContextOutput) Body(content []byte) {
	writer := this.Context.ResponseWriter.(io.Writer)
	this.Header("Content-Length", strconv.Itoa(len(content)))
	if this.Status != 0 {
		this.Context.ResponseWriter.WriteHeader(this.Status)
		this.Status = 0
	}

	writer.Write(content)
	switch writer.(type) {
	case *flate.Writer:
		writer.(*flate.Writer).Close()
	case *gzip.Writer:
		writer.(*gzip.Writer).Close()
	}
}

/**
 * Cookie
 * @param
 * 		key
 * 		value
 * 		maxAge,path,domain,secure,httponly
 * @return
 */
func (this *ContextOutput) Cookie(key, value string, params ...interface{}) {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "%s=%s", cookieKeyFormatter.Replace(key), cookieValueFormatter.Replace(value))

	//max age
	if len(params) > 0 {
		switch v := params[0].(type) {
		case int:
			if v > 0 {
				fmt.Fprintf(&buffer, "; Expires=%s; Maxage=%d", time.Now().Add(time.Duration(v)*time.Second).UTC().Format(time.RFC1123), v)
			} else {
				fmt.Fprintf(&buffer, "; Max-Age=0")
			}
		case int32:
			if v > 0 {
				fmt.Fprintf(&buffer, "; Expires=%s; Maxage=%d", time.Now().Add(time.Duration(v)*time.Second).UTC().Format(time.RFC1123), v)
			} else {
				fmt.Fprintf(&buffer, "; Max-Age=0")
			}
		case int64:
			if v > 0 {
				fmt.Fprintf(&buffer, "; Expires=%s; Maxage=%d", time.Now().Add(time.Duration(v)*time.Second).UTC().Format(time.RFC1123), v)
			} else {
				fmt.Fprintf(&buffer, "; Max-Age=0")
			}
		}
	}

	//path
	if len(params) > 1 {
		if v, ok := params[1].(string); ok {
			fmt.Fprintf(&buffer, "; Path=%s", cookieValueFormatter.Replace(v))
		} else {
			fmt.Fprintf(&buffer, "; Path=%s", "/")
		}
	}

	//domain
	if len(params) > 2 {
		if v, ok := params[2].(string); ok {
			fmt.Fprintf(&buffer, "; Domain=%s", cookieValueFormatter.Replace(v))
		}
	}

	//secure
	if len(params) > 3 {
		var secure bool
		switch params[3].(type) {
		case bool:
			secure = params[3].(bool)
		default:
			if params[3] != nil {
				secure = true
			}
		}
		if secure {
			fmt.Fprintf(&buffer, " ; Secure")
		}
	}

	//http only
	if len(params) > 4 {
		httponly := false
		if v, ok := params[4].(bool); ok && v {
			httponly = true
		}
		if httponly {
			fmt.Fprintf(&buffer, " ; HttpOnly")
		}
	}

	this.Context.ResponseWriter.Header().Add("Set-Cookie", buffer.String())
}

var cookieKeyFormatter = strings.NewReplacer("\n", "-", "\r", "-")
var cookieValueFormatter = strings.NewReplacer("\n", " ", "\r", " ", ";", " ")

//serve json
func (this *ContextOutput) Json(content interface{}, hasIntent bool) error {
	this.Header("Content-Type", "application/json; charset=utf-8")
	var jsonBody []byte
	var err error
	if hasIntent {
		jsonBody, err = json.MarshalIndent(content, "", " ")
	} else {
		jsonBody, err = json.Marshal(content)
	}
	if err != nil {
		http.Error(this.Context.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return err
	}
	this.Body(jsonBody)
	return nil
}

//download
func (this *ContextOutput) Download(filepath string, filename ...string) {
	if len(filepath) > 0 && filepath[0] == '/' {
		filepath = filepath[1:]
	}
	this.Header("Content-Description", "File Transfer")
	this.Header("Content-Type", "application/octet-stream")
	if len(filename) > 0 && filename[0] != "" {
		this.Header("Content-Disposition", "attachment; filename="+filename[0])
	} else {
		this.Header("Content-Disposition", "attachment; filename="+path.Base(filepath))
	}
	this.Header("Content-Transfer-Encoding", "binary")
	this.Header("Expires", "0")
	this.Header("Cache-Control", "must-revalidate")
	this.Header("Pragma", "public")
	http.ServeFile(this.Context.ResponseWriter, this.Context.Request, filepath)
}

//static file
func (this *ContextOutput) StaticFile(filepath string) {
	if len(filepath) > 0 && filepath[0] == '/' {
		filepath = filepath[1:]
	}
	http.ServeFile(this.Context.ResponseWriter, this.Context.Request, filepath)
}

//set status
func (this *ContextOutput) SetStatus(status int) {
	this.Status = status
}
