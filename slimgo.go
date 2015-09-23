package slimgo

import (
	"runtime"
	"strconv"
	"strings"
)

const VERSION = "1.02"

func Run(address ...string) {
	//多核
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(address) > 0 && address[0] != "" {
		spAddr := strings.Split(address[0], ":")
		if len(spAddr) > 0 && spAddr[0] != "" {
			HttpAddress = spAddr[0]
		}
		if len(spAddr) > 1 && spAddr[1] != "" {
			HttpPort, _ = strconv.Atoi(spAddr[1])
		}
	}
	SlimApp.Run()
}
