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

func (this *App) Run() {
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
