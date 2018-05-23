package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mickep76/qry"

	"github.com/imc-trading/dock2box/model"
)

func (h *Handler) AllHosts(w http.ResponseWriter, r *http.Request) {
	q, err := qry.FromURL(r.URL.Query())
	if err != nil {
		writeError(w, err)
		return
	}

	hosts, err := h.ds.QueryHosts(q)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, hosts)
}

func (h *Handler) OneHost(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	host, err := h.ds.OneHost(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, host)
}

func (h *Handler) CreateHost(w http.ResponseWriter, r *http.Request) {
	host := &model.Host{}
	if err := json.NewDecoder(r.Body).Decode(host); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.CreateHost(host); err != nil {
		writeError(w, err)
		return
	}

	write(w, host)
}

func (h *Handler) UpdateHost(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	host, err := h.ds.OneHost(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(host); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.UpdateHost(host); err != nil {
		writeError(w, err)
		return
	}

	write(w, host)
}
