package slimgo

import (
	"fmt"
	//"net"
	"net/http"
	"time"
)

type App struct {
	Server     *http.Server
	Handerlers *ControllerRegister
}

func NewApp() *App {
	ctrlReg := NewControllerRegister()
	app := &App{
		Server:     &http.Server{},
		Handerlers: ctrlReg,
	}
	return app
}

func (this *App) beforeRun() {
	//timetask
	for name, t := range tasks {
		if t.runtag == true {
			t.Start()
			Logger.Info("Timetask:", name, " is running.")
		}
	}

	//hooks
	for key, f := range hooks.hookBeforeAppRun {
		err := f()
		if err != nil {
			Logger.Error("Hook:", key, " err,", err.Error())
		} else {
			Logger.Info("Hook:", key, " finished.")
		}
	}
}

func (this *App) Run() {
	this.beforeRun()
	address := HttpAddress
	if HttpPort != 0 {
		address = fmt.Sprintf("%s:%d", address, HttpPort)
	}
	var (
		err error
		//listener net.Listener
	)
	ending := make(chan bool, 1)
	//http
	this.Server.Addr = address
	this.Server.Handler = this.Handerlers
	this.Server.ReadTimeout = time.Duration(HttpServerTimeOut) * time.Second
	this.Server.WriteTimeout = time.Duration(HttpServerTimeOut) * time.Second
	Logger.Info("Http running on ", this.Server.Addr)
	go func() {
		err = this.Server.ListenAndServe()
		if err != nil {
			Logger.Error(err.Error())
			time.Sleep(100 * time.Microsecond)
			ending <- true
		}
	}()
	<-ending
}
