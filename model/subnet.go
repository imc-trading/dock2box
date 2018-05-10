package model

import (
	"fmt"
	"time"

	"github.com/mickep76/qry"
	"github.com/pborman/uuid"
)

type Subnet struct {
	UUID       string     `json:"uuid"`
	Created    time.Time  `json:"created"`
	Updated    *time.Time `json:"updated,omitempty"`
	Name       string     `json:"name"`
	Network    string     `json:"network"`
	CIDR       int        `json:"cidr"`
	Gateway    string     `json:"gateway,omitempty"`
	Site       *Site      `json:"site,omitempty"`
	DNSServers []string   `json:"dnsServers,omitempty"`
	DNSSearch  []string   `json:"dnsSearch,omitempty"`
}

type Subnets []*Subnet

func NewSubnet(name string) *Subnet {
	return &Subnet{
		UUID:    uuid.New(),
		Created: time.Now(),
		Name:    name,
	}
}

func (ds *Datastore) QuerySubnet(q *qry.Query) (Subnets, error) {
	kvs, err := ds.Values("subnets")
	if err != nil {
		return nil, err
	}

	subnets := Subnets{}
	if err := kvs.Decode(&subnets); err != nil {
		return nil, err
	}

	r, err := q.Eval(subnets)
	if err != nil {
		return nil, err
	}

	return r.(Subnets), nil
}

func (ds *Datastore) OneSubnet(uuid string) (*Subnet, error) {
	kvs, err := ds.Values(fmt.Sprintf("subnets/%s", uuid))
	if err != nil {
		return nil, err
	}

	subnets := Subnets{}
	if err := kvs.Decode(&subnets); err != nil {
		return nil, err
	}

	if len(subnets) > 0 {
		return subnets[0], nil
	}

	return nil, nil
}

func (ds *Datastore) CreateSubnet(subnet *Subnet) error {
	return ds.Set(fmt.Sprintf("subnets/%s", subnet.UUID), subnet)
}

func (ds *Datastore) UpdateSubnet(subnet *Subnet) error {
	now := time.Now()
	subnet.Updated = &now
	return ds.Set(fmt.Sprintf("subnets/%s", subnet.UUID), subnet)
}
