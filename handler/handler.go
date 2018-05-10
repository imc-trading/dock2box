package handler

import (
	"github.com/imc-trading/dock2box/model"
)

type Handler struct {
	ds *model.Datastore
}

func NewHandler(ds *model.Datastore) *Handler {
	return &Handler{
		ds: ds,
	}
}
