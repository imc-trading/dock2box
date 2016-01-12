package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

// LoadFiles returns field from multiple files.
func LoadFiles(files []string) (map[string]string, error) {
	r := make(map[string]string)

	for _, fn := range files {
		if _, err := os.Stat(fn); os.IsNotExist(err) {
			return map[string]string{}, fmt.Errorf("file doesn't exist: %s", fn)
		}

		o, err := ioutil.ReadFile(fn)
		if err != nil {
			return map[string]string{}, err
		}

		r[path.Base(fn)] = strings.TrimRight(string(o), "\n")
	}

	return r, nil
}

// LoadFileFields returns fields from file.
func LoadFileFields(fn string, del string, fields []string) (map[string]string, error) {
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		return map[string]string{}, fmt.Errorf("file doesn't exist: %s", fn)
	}

	o, err := ioutil.ReadFile(fn)
	if err != nil {
		return map[string]string{}, err
	}

	r, err := parseFields(string(o), del, fields)
	if err != nil {
		return map[string]string{}, err
	}

	return r, nil
}

// ExecCmd returns output.
func ExecCmd(cmd string, args []string) (string, error) {
	_, err := exec.LookPath(cmd)
	if err != nil {
		return "", fmt.Errorf("command doesn't exist: %s", cmd)
	}

	o, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimRight(string(o), "\n"), err
}

// ExecCmdFields returns fields from output.
func ExecCmdFields(cmd string, args []string, del string, fields []string) (map[string]string, error) {
	_, err := exec.LookPath(cmd)
	if err != nil {
		return map[string]string{}, fmt.Errorf("command doesn't exist: %s", cmd)
	}

	o, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return map[string]string{}, err
	}

	r, err := parseFields(string(o), del, fields)
	if err != nil {
		return map[string]string{}, err
	}

	return r, nil
}

func parseFields(o string, del string, fields []string) (map[string]string, error) {
	r := make(map[string]string)

	for _, line := range strings.Split(string(o), "\n") {
		vals := strings.SplitN(line, del, 2)
		if len(vals) < 1 {
			continue
		}

		for _, f := range fields {
			if strings.HasPrefix(strings.TrimLeft(line, " \t"), f) {
				r[f] = strings.Trim(strings.Join(vals[1:], " "), " \t")
			}
		}
	}

	return r, nil
}
