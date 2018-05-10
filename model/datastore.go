package model

import (
	"github.com/mickep76/kvstore"
	_ "github.com/mickep76/kvstore/etcdv3"
)

type Datastore struct {
	lease kvstore.Lease
	kvstore.Conn
}

func NewDatastore(driver string, endpoints []string, keepalive int, options ...func(kvstore.Driver) error) (*Datastore, error) {
	c, err := kvstore.Open(driver, endpoints, options...)
	if err != nil {
		return nil, err
	}

	l, err := c.Lease(keepalive)
	if err != nil {
		return nil, err
	}

	return &Datastore{
		lease: l,
		Conn:  c,
	}, nil
}

func (ds *Datastore) Lease() kvstore.Lease {
	return ds.lease
}
