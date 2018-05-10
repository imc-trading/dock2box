package model

import "time"

type Rack struct {
	UUID       string     `json:"uuid"`
	Created    time.Time  `json:"created"`
	Updated    *time.Time `json:"updated,omitempty"`
	Name       string     `json:"name"`
	Descr      string     `json:"descr"`
	Site       *Site      `json:"site"`
	UnitStart  int        `json:"unitStart"`
	UnitHeight int        `json:"unitHeight"`
}

type Racks []*Rack
