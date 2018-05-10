package model

import "time"

type HostGroup struct {
	UUID    string     `json:"uuid"`
	Created time.Time  `json:"created"`
	Updated *time.Time `json:"updated,omitempty"`
	Name    string     `json:"name"`
	Descr   string     `json:"descr"`
}

type HostGroups []*HostGroup
