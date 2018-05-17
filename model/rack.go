package model

import (
	"fmt"
	"time"

	"github.com/mickep76/qry"
	"github.com/pborman/uuid"
)

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

func NewRack(name string) *Rack {
	return &Rack{
		UUID:    uuid.New(),
		Created: time.Now(),
		Name:    name,
	}
}

func (ds *Datastore) AllRacks() (Racks, error) {
	kvs, err := ds.Values("racks")
	if err != nil {
		return nil, err
	}

	racks := Racks{}
	if err := kvs.Decode(&racks); err != nil {
		return nil, err
	}

	return racks, nil
}

func (ds *Datastore) QueryRacks(q *qry.Query) (Racks, error) {
	racks, err := ds.AllRacks()
	if err != nil {
		return nil, err
	}

	filtered, err := q.Query(racks)
	if err != nil {
		return nil, err
	}

	return filtered.(Racks), nil
}

func (ds *Datastore) OneRack(uuid string) (*Rack, error) {
	kvs, err := ds.Values(fmt.Sprintf("racks/%s", uuid))
	if err != nil {
		return nil, err
	}

	racks := Racks{}
	if err := kvs.Decode(&racks); err != nil {
		return nil, err
	}

	if len(racks) > 0 {
		return racks[0], nil
	}

	return nil, nil
}

func (ds *Datastore) CreateRack(rack *Rack) error {
	return ds.Set(fmt.Sprintf("racks/%s", rack.UUID), rack)
}

func (ds *Datastore) UpdateRack(rack *Rack) error {
	now := time.Now()
	rack.Updated = &now
	return ds.Set(fmt.Sprintf("racks/%s", rack.UUID), rack)
}
