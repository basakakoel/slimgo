package slimgo

import (
	"github.com/Unknwon/goconfig"
	"os"
	"path/filepath"
	"strings"
)

var (
	SlimApp           *App
	AppName           string
	WorkPath          string
	AppPath           string
	HttpAddress       string
	HttpPort          int
	HttpServerTimeOut int64
	SessionOn         bool
	StaticPath        map[string]string
)

func init() {
	SlimApp = NewApp()

	AppName = "SlimGo"

	WorkPath, _ = os.Getwd()
	WorkPath, _ = filepath.Abs(WorkPath)
	AppPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	HttpAddress = ""
	HttpPort = 6969

	HttpServerTimeOut = 0

	SessionOn = false

	StaticPath = map[string]string{
		"/public": "public",
	}

	CoverConfigByUser()
}

//覆盖配置
func CoverConfigByUser() {
	confPath := "conf/conf.ini"
	coverConfs, err := goconfig.LoadConfigFile(confPath)
	if err != nil {
		Waring("Can't find config file:conf/conf.ini ! Use default config.", err.Error())
		return
	}

	if appname, err := coverConfs.GetValue(goconfig.DEFAULT_SECTION, "AppName"); err == nil {
		AppName = appname
	}

	if httpAddress, err := coverConfs.GetValue(goconfig.DEFAULT_SECTION, "HttpAddress"); err == nil {
		HttpAddress = httpAddress
	}

	if httpPort, err := coverConfs.Int(goconfig.DEFAULT_SECTION, "HttpPort"); err == nil {
		HttpPort = httpPort
	}

	if httpServerTimeOut, err := coverConfs.Int64(goconfig.DEFAULT_SECTION, "HttpServerTimeOut"); err == nil {
		HttpServerTimeOut = httpServerTimeOut
	}

	if sessionOn, err := coverConfs.Bool(goconfig.DEFAULT_SECTION, "SessionOn"); err == nil {
		SessionOn = sessionOn
	}
}

//set static path
func SetStaticPath(url, path string) {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	url = strings.TrimRight(url, "/")
	StaticPath[url] = path
}

//delete static path
func DeleteStaticPath(url string) {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	url = strings.TrimRight(url, "/")
	delete(StaticPath, url)
}
