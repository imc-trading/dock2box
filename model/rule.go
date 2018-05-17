package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mickep76/qry"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
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

func (ds *Datastore) AllRules(q *qry.Query) (Rules, error) {
	return rules, nil
}

func (ds *Datastore) QueryRules(q *qry.Query) (Rules, error) {
	r, err := q.Query(rules)
	if err != nil {
		return nil, err
	}

	return r.(Rules), nil
}

func RunRules(h *Host, rules Rules) (map[string]string, error) {
	vm := otto.New()

	res := map[string]string{}
	for _, r := range rules {
		if r.Code != "" {
			// Convert to map[string]interface{}
			var m map[string]interface{}
			b, _ := json.Marshal(h)
			json.Unmarshal(b, &m)

			vm.Set("host", m)
			v, err := vm.Run(r.Code)
			if err != nil {
				return nil, errors.Wrapf(err, "run rule: %s host: %s", r.Name, h.Name)
			}

			s := fmt.Sprintf("%v", v)
			if s != "undefined" {
				res[r.Name] = s
			}
		} else if r.Value != "" {
			res[r.Name] = fmt.Sprintf("%v", r.Value)
		}
	}

	return res, nil
}
