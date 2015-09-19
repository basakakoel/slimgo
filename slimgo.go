package slimgo

import (
	"strconv"
	"strings"
)

func Run(address ...string) {
	if len(address) > 0 && address[0] != "" {
		spAddr := strings.Split(address[0], ":")
		if len(spAddr) > 0 && spAddr[0] != "" {
			HttpAddres = spAddr[0]
		}
		if len(spAddr) > 1 && spAddr[1] != "" {
			HttpPort, _ = strconv.Atoi(spAddr[1])
		}
	}
	SlimApp.Run()
}
