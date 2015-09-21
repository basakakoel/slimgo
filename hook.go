package slimgo

import (
	"github.com/jesusslim/slimgo/context"
)

type Hooks struct {
	hookBeforeAppRun    map[string]func() error
	hookBeforeHttpPre   map[string]func(*context.Context)
	hookAfterHttpFinish map[string]func(*context.Context)
}

var hooks = &Hooks{
	hookBeforeAppRun:    make(map[string]func() error),
	hookBeforeHttpPre:   make(map[string]func(*context.Context)),
	hookAfterHttpFinish: make(map[string]func(*context.Context)),
}

func RegisterHook(where, key string, theFunc interface{}) {
	switch where {
	case "BAR":
		ResisterHookBeforeAppRun(key, theFunc.(func() error))
		break
	case "BHP":
		RegisterHookBeforeHttpPre(key, theFunc.(func(*context.Context)))
		break
	case "AHF":
		RegisterHookAfterHttpFinish(key, theFunc.(func(*context.Context)))
		break
	}
}

func ResisterHookBeforeAppRun(key string, theFunc func() error) {
	hooks.hookBeforeAppRun[key] = theFunc
}

func RegisterHookBeforeHttpPre(key string, theFunc func(*context.Context)) {
	hooks.hookBeforeHttpPre[key] = theFunc
}

func RegisterHookAfterHttpFinish(key string, theFunc func(*context.Context)) {
	hooks.hookAfterHttpFinish[key] = theFunc
}
