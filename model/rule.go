package model

import (
	"time"

	"github.com/mickep76/qry"
	"github.com/pborman/uuid"
)

type Rule struct {
	UUID    string      `json:"uuid"`
	Created time.Time   `json:"created"`
	Updated *time.Time  `json:"updated,omitempty"`
	Name    string      `json:"name" toml:"name"`
	Descr   string      `json:"descr" toml:"descr"`
	Value   interface{} `json:"value,omitempty" toml:"value"`
	Code    string      `json:"code,omitempty" toml:"code"`
	File    string      `json:"file,omitempty" toml:"file"`
}

type Rules []*Rule

var rules Rules

func NewRule(name string, descr string, value interface{}, code string, file string) *Rule {
	return &Rule{
		UUID:    uuid.New(),
		Created: time.Now(),
		Name:    name,
		Descr:   descr,
		Value:   value,
		Code:    code,
		File:    file,
	}
}

func (ds *Datastore) QueryRules(q *qry.Query) (Rules, error) {
	r, err := q.Query(rules)
	if err != nil {
		return nil, err
	}

	return r.(Rules), nil
}
