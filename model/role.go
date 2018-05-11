package model

import (
	"fmt"
	"time"

	"github.com/mickep76/qry"
	"github.com/pborman/uuid"
)

type Role struct {
	UUID    string     `json:"uuid"`
	Created time.Time  `json:"created"`
	Updated *time.Time `json:"updated,omitempty"`
	Name    string     `json:"name"`
}

type Roles []*Role

func NewRole(name string) *Role {
	return &Role{
		UUID:    uuid.New(),
		Created: time.Now(),
		Name:    name,
	}
}

func (ds *Datastore) QueryRoles(q *qry.Query) (Roles, error) {
	kvs, err := ds.Values("roles")
	if err != nil {
		return nil, err
	}

	roles := Roles{}
	if err := kvs.Decode(&roles); err != nil {
		return nil, err
	}

	r, err := q.Eval(roles)
	if err != nil {
		return nil, err
	}

	return r.(Roles), nil
}

func (ds *Datastore) OneRole(uuid string) (*Role, error) {
	kvs, err := ds.Values(fmt.Sprintf("roles/%s", uuid))
	if err != nil {
		return nil, err
	}

	roles := Roles{}
	if err := kvs.Decode(&roles); err != nil {
		return nil, err
	}

	if len(roles) > 0 {
		return roles[0], nil
	}

	return nil, nil
}

func (ds *Datastore) CreateRole(role *Role) error {
	return ds.Set(fmt.Sprintf("roles/%s", role.UUID), role)
}

func (ds *Datastore) UpdateRole(role *Role) error {
	now := time.Now()
	role.Updated = &now
	return ds.Set(fmt.Sprintf("roles/%s", role.UUID), role)
}
