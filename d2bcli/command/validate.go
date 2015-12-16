package command

import (
	"fmt"
	"net"
	"strings"
)

func validateHwAddr(inp string, dmy string) bool {
	if _, err := net.ParseMAC(inp); err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func validateIPv4(inp string, dmy string) bool {
	if net.ParseIP(inp) == nil {
		fmt.Println("Invalid IPv4 address")
		return false
	}
	return true
}

func validateIPv4List(inp string, list string) bool {
	for _, v := range strings.Split(inp, ",") {
		if net.ParseIP(v) == nil {
			fmt.Println("Invalid IPv4 address: %s", v)
			return false
		}
	}
	return true
}
