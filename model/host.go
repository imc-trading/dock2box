package model

import (
	"fmt"
	"time"

	"github.com/mickep76/qry"
	"github.com/pborman/uuid"
)

type Host struct {
	UUID           string         `json:"uuid"`
	Created        time.Time      `json:"created"`
	Updated        *time.Time     `json:"updated,omitempty"`
	ClientUUID     string         `json:"clientUUID"`
	Name           string         `json:"name"`
	IPAddr         string         `json:"ipAddr"`
	Hardware       *Hardware      `json:"hardware,omitempty"`
	HwAddr         string         `json:"hwAddr,omitempty"`
	SerialNumber   string         `json:"serialNumber,omitempty"`
	Role           *Role          `json:"role,omitempty"`
	Host           *Host          `json:"host,omitempty"`
	Site           *Site          `json:"site,omitempty"`
	Rack           *Rack          `json:"rack,omitempty"`
	RackUnitPos    int            `json:"rackUnitPos,omitempty"`
	RackUnitHeight int            `json:"rackUnitHeight,omitempty"`
	Tenant         *Tenant        `json:"tenantUUID,omitempty"`
	DockerImage    *DockerImage   `json:"dockerImage,omitempty"`
	AllowBuild     bool           `json:"allowBuild"`
	SecurityGroup  *SecurityGroup `json:"securityGroup"`
	HostGroups     HostGroups     `json:"hostGroups,omitempty"`
}

type Hosts []*Host

func NewHost(name string, ipAddr string) *Host {
	return &Host{
		UUID:    uuid.New(),
		Created: time.Now(),
		Name:    name,
		IPAddr:  ipAddr,
	}
}

func (ds *Datastore) QueryHost(q *qry.Query) (Hosts, error) {
	kvs, err := ds.Values("hosts")
	if err != nil {
		return nil, err
	}

	hosts := Hosts{}
	if err := kvs.Decode(&hosts); err != nil {
		return nil, err
	}

	r, err := q.Eval(hosts)
	if err != nil {
		return nil, err
	}

	return r.(Hosts), nil
}

func (ds *Datastore) OneHost(uuid string) (*Host, error) {
	kvs, err := ds.Values(fmt.Sprintf("hosts/%s", uuid))
	if err != nil {
		return nil, err
	}

	hosts := Hosts{}
	if err := kvs.Decode(&hosts); err != nil {
		return nil, err
	}

	if len(hosts) > 0 {
		return hosts[0], nil
	}

	return nil, nil
}

func (ds *Datastore) CreateHost(host *Host) error {
	return ds.Set(fmt.Sprintf("hosts/%s", host.UUID), host)
}

func (ds *Datastore) UpdateHost(host *Host) error {
	now := time.Now()
	host.Updated = &now
	return ds.Set(fmt.Sprintf("hosts/%s", host.UUID), host)
}
