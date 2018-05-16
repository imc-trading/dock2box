package model

import (
	"fmt"
	"time"

	"github.com/pborman/uuid"
)

type Site struct {
	UUID    string     `json:"uuid"`
	Created time.Time  `json:"created"`
	Updated *time.Time `json:"updated,omitempty"`
	Name    string     `json:"name"`
	Descr   string     `json:"descr,omitempty"`
	Region  string     `json:"region,omitempty"`
	Country string     `json:"country,omitempty"`
	Lat     float32    `json:"lat,omitempty"`
	Lng     float32    `json:"lng,omitempty"`
}

type Sites []*Site

func NewSite(name string) *Site {
	return &Site{
		UUID:    uuid.New(),
		Created: time.Now(),
		Name:    name,
	}
}

func (ds *Datastore) AllSites() (Sites, error) {
	kvs, err := ds.Values("sites")
	if err != nil {
		return nil, err
	}

	sites := Sites{}
	if err := kvs.Decode(&sites); err != nil {
		return nil, err
	}

	return sites, nil
}

func (ds *Datastore) OneSite(uuid string) (*Site, error) {
	kvs, err := ds.Values(fmt.Sprintf("sites/%s", uuid))
	if err != nil {
		return nil, err
	}

	sites := Sites{}
	if err := kvs.Decode(&sites); err != nil {
		return nil, err
	}

	if len(sites) > 0 {
		return sites[0], nil
	}

	return nil, nil
}

func (ds *Datastore) CreateSite(site *Site) error {
	return ds.Set(fmt.Sprintf("sites/%s", site.UUID), site)
}

func (ds *Datastore) UpdateSite(site *Site) error {
	now := time.Now()
	site.Updated = &now
	return ds.Set(fmt.Sprintf("sites/%s", site.UUID), site)
}
