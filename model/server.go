package model

import (
	"fmt"
	"time"

	"github.com/mickep76/kvstore"
	"github.com/pborman/uuid"
)

type Server struct {
	UUID    string     `json:"uuid"`
	Created time.Time  `json:"created"`
	Updated *time.Time `json:"updated,omitempty"`
	Name    string     `json:"name"`
}

type Servers []*Server

func NewServer(name string) *Server {
	return &Server{
		UUID:    uuid.New(),
		Created: time.Now(),
		Name:    name,
	}
}

func (ds *Datastore) AllServers() (Servers, error) {
	kvs, err := ds.Values("servers")
	if err != nil {
		return nil, err
	}

	servers := Servers{}
	if err := kvs.Decode(&servers); err != nil {
		return nil, err
	}

	return servers, nil
}

func (ds *Datastore) OneServer(uuid string) (*Server, error) {
	kvs, err := ds.Values(fmt.Sprintf("servers/%s", uuid))
	if err != nil {
		return nil, err
	}

	servers := Servers{}
	if err := kvs.Decode(&servers); err != nil {
		return nil, err
	}

	if len(servers) > 0 {
		return servers[0], nil
	}

	return nil, nil
}

func (ds *Datastore) CreateServer(server *Server) error {
	return ds.Set(fmt.Sprintf("servers/%s", server.UUID), server, kvstore.WithLease(ds.lease))
}

func (ds *Datastore) UpdateServer(server *Server) error {
	now := time.Now()
	server.Updated = &now
	return ds.Set(fmt.Sprintf("servers/%s", server.UUID), server, kvstore.WithLease(ds.lease))
}
