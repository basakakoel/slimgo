package slimgo

import (
	"time"
)

type TimeTask struct {
	theFunc  func() error
	duration int64
	runtag   bool
}

func (this *TimeTask) Start() {

}

func (this *TimeTask) Run() {
	if this.runtag {
		err := this.theFunc()
		if err != nil {
			Logger.Error(err.Error())
		}
	}
}
