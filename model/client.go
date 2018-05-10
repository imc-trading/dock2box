package model

import (
	"fmt"
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
	IPAddr   string     `json:"ipAddr"`
	Hardware *Hardware  `json:"hardware,omitempty"`
}

type Clients []*Client

func NewClient(name string) *Client {
	return &Client{
		UUID:    uuid.New(),
		Created: time.Now(),
		Name:    name,
	}
}

func (ds *Datastore) QueryClients(q *qry.Query) (Clients, error) {
	kvs, err := ds.Values("clients")
	if err != nil {
		return nil, err
	}

	clients := Clients{}
	if err := kvs.Decode(&clients); err != nil {
		return nil, err
	}

	r, err := q.Eval(clients)
	if err != nil {
		return nil, err
	}

	return r.(Clients), nil
}

func (ds *Datastore) CreateClient(client *Client) error {
	return ds.Set(fmt.Sprintf("clients/%s", client.UUID), client, kvstore.WithLease(ds.lease))
}

func (ds *Datastore) UpdateClient(client *Client) error {
	now := time.Now()
	client.Updated = &now
	return ds.Set(fmt.Sprintf("clients/%s", client.UUID), client, kvstore.WithLease(ds.lease))
}
