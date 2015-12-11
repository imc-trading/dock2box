// +build linux

package routes

import (
	"os/exec"
	"strconv"
	"strings"
)

func Get() ([]Route, error) {
	routes := []Route{}

	o, err := exec.Command("netstat", "-rn").Output()
	if err != nil {
		return []Route{}, err
	}

	for c, line := range strings.Split(string(o), "\n") {
		vals := strings.Fields(line)
		if c < 2 || len(vals) < 8 {
			continue
		}

		r := Route{}

		r.Destination = vals[0]
		r.Gateway = vals[1]
		r.Genmask = vals[2]
		r.Flags = vals[3]

		r.MSS, err = strconv.Atoi(vals[4])
		if err != nil {
			return []Route{}, err
		}

		r.Window, err = strconv.Atoi(vals[5])
		if err != nil {
			return []Route{}, err
		}

		r.IRTT, err = strconv.Atoi(vals[6])
		if err != nil {
			return []Route{}, err
		}

		r.Interface = vals[7]

		routes = append(routes, r)
	}

	return routes, nil
}
