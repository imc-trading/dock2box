package model

import "time"

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
