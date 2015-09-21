package slimgo

import (
	"errors"
	"time"
)

var tasks map[string]*TimeTask

type TimeTask struct {
	theFunc  func() error
	duration int
	runtag   bool
}

func init() {
	tasks = make(map[string]*TimeTask)
}

func (this *TimeTask) Start() {
	this.Run()
}

func (this *TimeTask) Run() {
	if this.runtag {
		err := this.theFunc()
		if err != nil {
			Logger.Error(err.Error())
		}
		time.AfterFunc(time.Duration(this.duration)*time.Second, func() { this.Run() })
	}
}

func RegisterTimeTask(name string, theFunc func() error, duration int, runtag bool) {
	timeTask := &TimeTask{
		theFunc:  theFunc,
		duration: duration,
		runtag:   runtag,
	}
	tasks[name] = timeTask
}

func ShutDownTimeTask(name string) error {
	task, ok := tasks[name]
	if !ok {
		return errors.New(name + " not found")
	}
	if task.runtag == false {
		return errors.New(name + " is not running")
	}
	task.runtag = false
	return nil
}

func RestartTimeTask(name string) error {
	task, ok := tasks[name]
	if !ok {
		return errors.New(name + " not found")
	}
	if task.runtag == true {
		return errors.New(name + " is running")
	}
	task.runtag = true
	task.Start()
	return nil
}

func IsRunningTimeTask(name string) (bool, error) {
	task, ok := tasks[name]
	if !ok {
		return false, errors.New(name + " not found")
	}
	return task.runtag, nil
}
