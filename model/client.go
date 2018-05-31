package model

import (
	"fmt"
	"log"
	"time"

	"github.com/mickep76/kvstore"
	"github.com/mickep76/qry"
	"github.com/pborman/uuid"
)

type Client struct {
	UUID     string     `json:"uuid"`
	Created  time.Time  `json:"created"`
	Updated  *time.Time `json:"updated,omitempty"`
	Name     string     `json:"name"`
	Hardware *Hardware  `json:"hardware,omitempty"`
	HostUUID string     `json:"hostUUID"`
	Host     *Host      `json:"host"`
}

type Clients []*Client

func NewClient(name string) *Client {
	hw, errs := NewHardware()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Print(err)
		}
	}

	return &Client{
		UUID:     uuid.New(),
		Created:  time.Now(),
		Name:     name,
		Hardware: hw,
	}
}

func (ds *Datastore) AllClients() (Clients, error) {
	kvs, err := ds.Values("clients")
	if err != nil {
		return nil, err
	}

	clients := Clients{}
	if err := kvs.Decode(&clients); err != nil {
		return nil, err
	}

	return clients, nil
}

func (ds *Datastore) QueryClients(q *qry.Query) (Clients, error) {
	clients, err := ds.AllClients()
	if err != nil {
		return nil, err
	}

	filtered, err := q.Query(clients)
	if err != nil {
		return nil, err
	}

	return filtered.(Clients), nil
}

func (ds *Datastore) OneClient(uuid string) (*Client, error) {
	kvs, err := ds.Values(fmt.Sprintf("clients/%s", uuid))
	if err != nil {
		return nil, err
	}

	clients := Clients{}
	if err := kvs.Decode(&clients); err != nil {
		return nil, err
	}

	if len(clients) > 0 {
		return clients[0], nil
	}

	return nil, nil
}

func (ds *Datastore) CreateClient(client *Client) error {
	return ds.Set(fmt.Sprintf("clients/%s", client.UUID), client, kvstore.WithLease(ds.lease))
}

func (ds *Datastore) UpdateClient(client *Client) error {
	now := time.Now()
	client.Updated = &now
	return ds.Set(fmt.Sprintf("clients/%s", client.UUID), client, kvstore.WithLease(ds.lease))
}
