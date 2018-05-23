package handler

import (
	"encoding/json"
	"net/http"

	"github.com/imc-trading/dock2box/model"

	"github.com/mickep76/auth"
	_ "github.com/mickep76/auth/ldap"
)

type Handler struct {
	ds   *model.Datastore
	jwt  *auth.JWT
	conn auth.Conn
}

func NewHandler(ds *model.Datastore) *Handler {
	return &Handler{
		ds: ds,
	}
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

func write(w http.ResponseWriter, v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
