package slimgo

import (
	"slimgo/context"
	"strconv"
)

func exception(errcode string, ctx *context.Context) {
	code, err := strconv.Atoi(errcode)
	if err != nil {
		code = 503
	}
	ctx.ResponseWriter.WriteHeader(code)
	ctx.WriteString(errcode)
}
