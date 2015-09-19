package slimgo

import (
	"slimgo/context"
)

type ControllerInterface interface {
	Init(ctx *context.Context, controllerName, actionName string)
	Pre()
	URLMapping()
	Finish()
	HandlerFunc(funcName string) bool
}

type Controller struct {
	Context        *context.Context
	Data           map[interface{}]interface{}
	controllerName string
	actionName     string
	funcMapping    map[string]func()
}

func (this *Controller) Init(ctx *context.Context, controllerName, actionName string) {
	this.Context = ctx
	this.Data = make(map[interface{}]interface{})
	this.controllerName = controllerName
	this.actionName = actionName
	this.funcMapping = make(map[string]func())
}

func (this *Controller) Pre() {

}

func (this *Controller) URLMapping() {

}

func (this *Controller) Finish() {

}

//call func
func (this *Controller) HandlerFunc(funcName string) bool {
	if theFunc, ok := this.funcMapping[funcName]; ok {
		theFunc()
		return true
	} else {
		return false
	}
}

/********** **********/

//mapping method to a map for router
func (this *Controller) Mapping(funcName string, theFunc func()) {
	this.funcMapping[funcName] = theFunc
}

func (this *Controller) Redirect(url string, status int) {
	this.Context.Redirect(url, status)
}

func (this *Controller) ServeJson(data interface{}) {
	hasIntent := true
	if data != nil {
		this.Context.Output.Json(data, hasIntent)
	} else {
		this.Context.Output.Json(this.Data["json"], hasIntent)
	}
}

//default index
func (this *Controller) Index() {
	this.Data["json"] = "Hello Slimgo."
	this.ServeJson(nil)
}
