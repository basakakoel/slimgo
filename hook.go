package slimgo

import (
	"github.com/jesusslim/slimgo/context"
)

type Hook struct {
	hookBeforeAppRun    map[string]func() error
	hookBeforeHttpPre   map[string]func(*context.Context)
	hookAfterHttpFinish map[string]func(*context.Context)
}

type TimeTask struct {
	theFunc  func() error
	duration int64
	runtag   bool
}
