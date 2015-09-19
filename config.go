package slimgo

import (
	"os"
	"path/filepath"
)

var (
	SlimApp           *App
	AppName           string
	WorkPath          string
	AppPath           string
	HttpAddres        string
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

	HttpAddres = ""
	HttpPort = 6969

	HttpServerTimeOut = 0

	SessionOn = false

	StaticPath = map[string]string{
		"/public": "public",
	}
}
