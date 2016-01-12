// +build linux

package sysctl

import (
	"os/exec"
	"strings"
)

func Get() ([]Sysctl, error) {
	sysctl := []Sysctl{}

	o, err := exec.Command("sysctl", "-a").Output()
	if err != nil {
		return []Sysctl{}, err
	}

	for _, line := range strings.Split(string(o), "\n") {
		vals := strings.Fields(line)
		if len(vals) < 3 {
			continue
		}

		s := Sysctl{}

		s.Key = vals[0]
		s.Value = vals[2]

		sysctl = append(sysctl, s)
	}

	return sysctl, nil
}
