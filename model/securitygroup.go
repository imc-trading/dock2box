package model

import "time"

type SecurityGroup struct {
	UUID    string     `json:"uuid"`
	Created time.Time  `json:"created"`
	Updated *time.Time `json:"updated,omitempty"`
	Name    string     `json:"name"`
}

type SecurityGroups []*SecurityGroup
