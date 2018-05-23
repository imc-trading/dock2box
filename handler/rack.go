package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mickep76/qry"

	"github.com/imc-trading/dock2box/model"
)

func (h *Handler) AllRacks(w http.ResponseWriter, r *http.Request) {
	q, err := qry.FromURL(r.URL.Query())
	if err != nil {
		writeError(w, err)
		return
	}

	racks, err := h.ds.QueryRacks(q)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, racks)
}

func (h *Handler) OneRack(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	rack, err := h.ds.OneHost(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, rack)
}

func (h *Handler) CreateRack(w http.ResponseWriter, r *http.Request) {
	rack := &model.Rack{}
	if err := json.NewDecoder(r.Body).Decode(rack); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.CreateRack(rack); err != nil {
		writeError(w, err)
		return
	}

	write(w, rack)
}

func (h *Handler) UpdateRack(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	rack, err := h.ds.OneRack(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(rack); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.UpdateRack(rack); err != nil {
		writeError(w, err)
		return
	}

	write(w, rack)
}
