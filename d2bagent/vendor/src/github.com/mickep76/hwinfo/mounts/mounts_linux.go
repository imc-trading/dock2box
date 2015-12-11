// +build linux

package mounts

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func Get() ([]Mount, error) {
	mounts := []Mount{}

	fn := "/proc/mounts"
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		return []Mount{}, fmt.Errorf("file doesn't exist: %s", fn)
	}

	o, err := ioutil.ReadFile(fn)
	if err != nil {
		return []Mount{}, err
	}

	for c, line := range strings.Split(string(o), "\n") {
		vals := strings.Fields(line)
		if c < 1 || len(vals) < 1 {
			continue
		}

		m := Mount{}

		m.Source = vals[0]
		m.Target = vals[1]
		m.FSType = vals[2]
		m.Options = vals[3]

		mounts = append(mounts, m)
	}

	return mounts, nil
}
