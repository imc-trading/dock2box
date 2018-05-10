package handler

import (
	"github.com/mickep76/kvstore/example/model"
)

type Handler struct {
	ds *model.Datastore
}

func NewHandler(ds *model.Datastore) *Handler {
	return &Handler{
		ds: ds,
	}
}
