package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mickep76/qry"

	"github.com/imc-trading/dock2box/model"
)

func (h *Handler) AllPools(w http.ResponseWriter, r *http.Request) {
	q, err := qry.FromURL(r.URL.Query())
	if err != nil {
		writeError(w, err)
		return
	}

	pools, err := h.ds.QueryPools(q)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, pools)
}

func (h *Handler) OnePool(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	pool, err := h.ds.OneHost(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, pool)
}

func (h *Handler) CreatePool(w http.ResponseWriter, r *http.Request) {
	pool := &model.Pool{}
	if err := json.NewDecoder(r.Body).Decode(pool); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.CreatePool(pool); err != nil {
		writeError(w, err)
		return
	}

	write(w, pool)
}

func (h *Handler) UpdatePool(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	pool, err := h.ds.OnePool(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(pool); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.UpdatePool(pool); err != nil {
		writeError(w, err)
		return
	}

	write(w, pool)
}
