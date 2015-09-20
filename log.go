package slimgo

import (
	"fmt"
)

var Logger *SlimLogger

type SlimLogger struct {
}

func (this *SlimLogger) baseLog(pre string, content ...interface{}) {
	fmt.Println("["+pre+"]:", content)
}

func (this *SlimLogger) Error(content ...interface{}) {
	this.baseLog("ERR", content)
}

func (this *SlimLogger) Warning(content ...interface{}) {
	this.baseLog("WARN", content)
}

func (this *SlimLogger) Info(content ...interface{}) {
	this.baseLog("INFO", content)
}

/***** be used for slimgo *****/

func Error(content ...interface{}) {
	Logger.Error(content)
}

func Waring(content ...interface{}) {
	Logger.Warning(content)
}

func Info(content ...interface{}) {
	Logger.Info(content)
}

func Log(content ...interface{}) {
	Info(content)
}
