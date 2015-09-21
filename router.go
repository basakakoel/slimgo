package slimgo

import (
	"github.com/jesusslim/slimgo/context"
	"github.com/jesusslim/slimgo/utils"
	"net/http"
	"path"
	"reflect"
	"strings"
	"time"
)

type ControllerRegister struct {
	routers map[string]*RouterTree
}

func NewControllerRegister() *ControllerRegister {
	return &ControllerRegister{
		routers: make(map[string]*RouterTree),
	}
}

type controllerInfo struct {
	pattern        string
	controllerType reflect.Type
	method         string
	routerType     int
}

const (
	routerTypeAuto = 1
)

var (
	//supported http methods
	HTTPMETHOD = map[string]string{
		"GET":     "GET",
		"POST":    "POST",
		"PUT":     "PUT",
		"DELETE":  "DELETE",
		"PATCH":   "PATCH",
		"OPTIONS": "OPTIONS",
		"HEAD":    "HEAD",
		"TRACE":   "TRACE",
		"CONNECT": "CONNECT",
	}
	//except from router tree method
	exceptMethod = []string{
		"Init",
		"Pre",
		"URLMapping",
		"Finish",
		"HandlerFunc",
		"Mapping",
		"Redirect",
		"ServeJson",
	}
)

//reg controller auto
func (this *ControllerRegister) Register(controllers ...ControllerInterface) {
	for _, c := range controllers {
		cVal := reflect.ValueOf(c)
		cType := cVal.Type()                       //指针type
		cRealType := reflect.Indirect(cVal).Type() //实际type
		// pkgName := cType.PkgPath()
		// pkgNameLower := strings.ToLower(pkgName)
		pkgPathLower := strings.ToLower(cRealType.PkgPath())
		index := strings.Index(pkgPathLower, "controller")
		if index <= 0 {
			panic("Controller class must be in controller packge!")
		}
		moduleName := ""
		moduleIndex := index + len("controller/")
		if moduleIndex < len(pkgPathLower) {
			moduleName = pkgPathLower[moduleIndex:]
		}
		controllerName := strings.ToLower(strings.TrimSuffix(cRealType.Name(), "Controller"))
		for i := 0; i < cType.NumMethod(); i++ {
			methodName := cType.Method(i).Name
			if utils.StringInSlice(methodName, exceptMethod) {
				continue
			}
			pattern := path.Join(moduleName, controllerName, strings.ToLower(methodName)) //小写
			node := &controllerInfo{
				pattern:        pattern,
				controllerType: cRealType,
				method:         methodName,
				routerType:     routerTypeAuto,
			}
			for _, v := range HTTPMETHOD {
				this.AddRouter(v, pattern, node)
			}
		}
	}
}

//add router
func (this *ControllerRegister) AddRouter(httpMethod, pattern string, node *controllerInfo) {
	if tree, ok := this.routers[httpMethod]; ok {
		tree.AddRouter(pattern, node)
	} else {
		tree := NewRouterTree()
		this.routers[httpMethod] = tree
		tree.AddRouter(pattern, node)
	}
}

//find router
func (this *ControllerRegister) FindRouter(httpMethod, url string) *controllerInfo {
	tree, ok := this.routers[httpMethod]
	if !ok {
		//not support method
		return nil
	}
	result := tree.FindRouter(url)
	if result == nil {
		return nil
	} else {
		return result.(*controllerInfo)
	}
}

//is statistic file
func isStaticFile(url string) (bool, string, bool) {
	url = path.Clean(url)
	var realpath string
	isRouter := true
	if url == "/favicon.ico" || url == "/robots.txt" {
		realpath = url[1:]
	} else {
		for k, v := range StaticPath {
			if len(k) == 0 {
				continue
			}
			index := strings.Index(url, k)
			if index >= 0 {
				isRouter = false
				indexSplit := index + len(k)
				if indexSplit < len(url) {
					realpath = v + url[indexSplit:]
					break
				}
			}
		}
	}
	if len(realpath) > 0 {
		isRouter = false
		if utils.FileExists(realpath) {
			return true, realpath, isRouter
		}
	}
	return false, "", isRouter
}

func (this *ControllerRegister) ServeHTTP(reponseWriter http.ResponseWriter, request *http.Request) {
	startTime := time.Now()
	url := request.URL.Path

	if ok, realpath, isRouter := isStaticFile(url); ok {
		http.ServeFile(reponseWriter, request, realpath)
		return
	} else if !isRouter {
		http.NotFound(reponseWriter, request)
		return
	}

	ctx := &context.Context{
		Request:        request,
		ResponseWriter: reponseWriter,
		Input:          context.NewContextInput(request),
		Output:         context.NewContextOutput(),
	}
	ctx.Output.Context = ctx

	ctrlInfo := this.FindRouter(request.Method, url)

	defer this.recoverPanic()

	if ctrlInfo == nil {
		exception("404", ctx)
		Logger.Error(url, " not match")
		goto AfterAll
	} else {
		//find router
		ctrlObj := reflect.New(ctrlInfo.controllerType)
		execController, ok := ctrlObj.Interface().(ControllerInterface)
		if !ok {
			panic("controller is not ControllerInterface")
		}
		execController.Init(ctx, ctrlInfo.controllerType.Name(), ctrlInfo.method)

		//hooks
		for _, f := range hooks.hookBeforeHttpPre {
			f(ctx)
		}

		execController.Pre()
		in := make([]reflect.Value, 0)
		method := ctrlObj.MethodByName(ctrlInfo.method)
		method.Call(in)
		execController.Finish()

		//hooks
		for _, f := range hooks.hookAfterHttpFinish {
			f(ctx)
		}

	}

AfterAll:
	duration := time.Since(startTime)
	Logger.Info("serve http:", url, ",take time:", duration.Seconds()*1000, "ms")
}

func (this *ControllerRegister) recoverPanic() {
	if err := recover(); err != nil {
		Logger.Error(err)
	}
}
