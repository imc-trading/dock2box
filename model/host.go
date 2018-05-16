package model

import (
	"fmt"
	"time"

	"github.com/pborman/uuid"
)

type Host struct {
	UUID              string         `json:"uuid"`
	Created           time.Time      `json:"created"`
	Updated           *time.Time     `json:"updated,omitempty"`
	Name              string         `json:"name"`
	HwAddr            string         `json:"hwAddr,omitempty"`
	Hardware          *Hardware      `json:"hardware,omitempty"`
	RoleUUID          string         `json:"roleUUID,omitempty"`
	Role              *Role          `json:"role,omitempty"`
	SiteUUID          string         `json:"siteUUID,omitempty"`
	Site              *Site          `json:"site,omitempty"`
	RackUUID          string         `json:"rackUUID,omitempty"`
	Rack              *Rack          `json:"rack,omitempty"`
	RackPos           int            `json:"rackPos,omitempty"`
	RackHeight        int            `json:"rackHeight,omitempty"`
	TenantUUID        string         `json:"tenantUUID,omitempty"`
	Tenant            *Tenant        `json:"tenant,omitempty"`
	ImageUUID         string         `json:"imageUUID,omitempty"`
	Image             *Image         `json:"image,omitempty"`
	AllowBuild        bool           `json:"allowBuild"`
	SecurityGroupUUID string         `json:"securityGroupUUID,omitempty"`
	SecurityGroup     *SecurityGroup `json:"securityGroup"`
	HostGroupUUIDs    []string       `json:"hostGroupUUIDs"`
	HostGroups        HostGroups     `json:"hostGroups,omitempty"`
}

type Hosts []*Host

func NewHost(name string, hwAddr string) *Host {
	return &Host{
		UUID:    uuid.New(),
		Created: time.Now(),
		Name:    name,
		HwAddr:  hwAddr,
	}
}

func (ds *Datastore) AllHosts() (Hosts, error) {
	kvs, err := ds.Values("hosts")
	if err != nil {
		return nil, err
	}

	hosts := Hosts{}
	if err := kvs.Decode(&hosts); err != nil {
		return nil, err
	}

	return hosts, nil
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
