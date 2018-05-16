package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/mickep76/encdec"
	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
)

type TaskDef struct {
	UUID               string     `json:"uuid"`
	Created            time.Time  `json:"created"`
	Updated            *time.Time `json:"updated,omitempty"`
	Name               string     `json:"name" toml:"name"`
	Descr              string     `json:"descr" toml:"descr"`
	User               string     `json:"user,omitempty" toml:"user"`
	Group              string     `json:"group,omitempty" toml:"group"`
	Dir                string     `json:"dir,omitempty" toml:"dir"`
	Cmd                string     `json:"cmd" toml:"cmd"`
	FileExist          string     `json:"fileExist,omitempty" toml:"fileExist"`
	FileNotExist       string     `json:"fileNotExist,omitempty" toml:"fileNotExist"`
	Reboot             bool       `json:"reboot" toml:"reboot"`
	RebootExitCode     *int       `json:"rebootExitCode,omitempty" toml:"rebootExitCode"`
	RebootFileExist    string     `json:"rebootFileExist,omitempty" toml:"rebootFileExist"`
	RebootFileNotExist string     `json:"rebootFileNotExist,omitempty" toml:"rebootFileNotExist"`
	Env                Rules      `json:"env,omitempty" toml:"env"`
	Args               Rules      `json:"args,omitempty" toml:"args"`
	Trigger            string     `json:"trigger,omitempty"`
	Require            string     `json:"require,omitempty"`
	Concurrency        *int       `json:"concurrency" toml:"concurrency"`
	Timeout            *int       `json:"timeout" toml:"timeout"`
	Pre                Rules      `json:"pre,omitempty" toml:"pre"`
	Post               Rules      `json:"post,omitempty" toml:"post"`
}

type TaskDefs []*TaskDef
type TaskDefMap map[string]*TaskDef

type HostTaskDef struct {
	EnvVals  map[string]string `json:"envVals"`
	ArgVals  []string          `json:"argVals"`
	PreVals  map[string]string `json:"preVals"`
	PostVals map[string]string `json:"postVals"`

	*TaskDef
}

type HostTaskDefs []*HostTaskDef

var (
	taskDefs   TaskDefs
	taskDefMap TaskDefMap
)

func LoadTaskDefs(path string) error {
	taskDefs = TaskDefs{}
	taskDefMap = make(TaskDefMap)

	matches, err := filepath.Glob(filepath.Join(path, "*", "taskdef.toml"))
	if err != nil {
		return errors.Wrapf(err, "glob files: %s", path)
	}

	for _, fn := range matches {
		td := &TaskDef{}
		if err := encdec.FromFile("toml", fn, td); err != nil {
			return errors.Wrapf(err, "read and decode: %s", fn)
		}

		dir := filepath.Dir(fn)
		fn := filepath.Join(dir, "trigger.js")
		if _, err := os.Stat(fn); !os.IsNotExist(err) {
			trigger, err := ioutil.ReadFile(fn)
			if err != nil {
				return errors.Wrapf(err, "read trigger: %s", fn)
			}
			td.Trigger = string(trigger)
		}

		fn = filepath.Join(dir, "require.js")
		if _, err := os.Stat(fn); !os.IsNotExist(err) {
			require, err := ioutil.ReadFile(fn)
			if err != nil {
				return errors.Wrapf(err, "read require: %s", fn)
			}
			td.Require = string(require)
		}

		for k, v := range td.Env {
			if v.File == "" {
				continue
			}

			fn = filepath.Join(dir, v.File)
			rule, err := ioutil.ReadFile(fn)
			if err != nil {
				return errors.Wrapf(err, "read env rule: %s", fn)
			}
			td.Env[k].Code = string(rule)
		}

		for k, v := range td.Args {
			if v.File == "" {
				continue
			}

			fn = filepath.Join(dir, v.File)
			rule, err := ioutil.ReadFile(fn)
			if err != nil {
				return errors.Wrapf(err, "read arg rule: %s", fn)
			}
			td.Args[k].Code = string(rule)
		}

		for k, v := range td.Pre {
			if v.File == "" {
				continue
			}

			fn = filepath.Join(dir, v.File)
			rule, err := ioutil.ReadFile(fn)
			if err != nil {
				return errors.Wrapf(err, "read pre rule: %s", fn)
			}
			td.Pre[k].Code = string(rule)
		}

		for k, v := range td.Post {
			if v.File == "" {
				continue
			}

			fn = filepath.Join(dir, v.File)
			rule, err := ioutil.ReadFile(fn)
			if err != nil {
				return errors.Wrapf(err, "read post rule: %s", fn)
			}
			td.Post[k].Code = string(rule)
		}

		if td.Concurrency == nil {
			c := 1
			td.Concurrency = &c
		}

		if td.Timeout == nil {
			t := 60 * 60
			td.Timeout = &t
		}

		taskDefMap[td.Name] = td
		taskDefs = append(taskDefs, td)
	}

	return nil
}

func AllTaskDefs() (TaskDefs, error) {
	return taskDefs, nil
}

func OneTaskDef(name string) (*TaskDef, error) {
	if v, ok := taskDefMap[name]; ok {
		return v, nil
	}
	return nil, nil
}

func TaskDefsByHost(h *Host) (HostTaskDefs, error) {
	vm := otto.New()

	// Set host info.
	var m map[string]interface{}
	b, _ := json.Marshal(h)
	json.Unmarshal(b, &m)
	vm.Set("host", m)

	list := HostTaskDefs{}
	for _, td := range taskDefs {
		// Run require.

		res, err := vm.Run(td.Require)
		if err != nil {
			return nil, fmt.Errorf("run require [%s]: %v", td.Name, err)
		}

		if !res.IsBoolean() {
			continue
		}

		ok, err := res.ToBoolean()
		if err != nil {
			return nil, fmt.Errorf("convert result to boolean: %v", err)
		}

		if !ok {
			continue
		}

		env, err := RunRules(h, td.Env)
		if err != nil {
			return nil, err
		}

		list = append(list, &HostTaskDef{
			EnvVals: env,
			TaskDef: td,
		})
	}

	return list, nil
}
