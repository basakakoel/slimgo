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
	this.baseLog("Error", content)
}

func (this *SlimLogger) Info(content ...interface{}) {
	this.baseLog("Info", content)
}
