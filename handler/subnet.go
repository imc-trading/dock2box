package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mickep76/qry"

	"github.com/imc-trading/dock2box/model"
)

func (h *Handler) AllSubnets(w http.ResponseWriter, r *http.Request) {
	q, err := qry.FromURL(r.URL.Query())
	if err != nil {
		writeError(w, err)
		return
	}

	subnets, err := h.ds.QuerySubnets(q)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, subnets)
}

func (h *Handler) OneSubnet(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	subnet, err := h.ds.OneHost(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, subnet)
}

func (h *Handler) CreateSubnet(w http.ResponseWriter, r *http.Request) {
	subnet := &model.Subnet{}
	if err := json.NewDecoder(r.Body).Decode(subnet); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.CreateSubnet(subnet); err != nil {
		writeError(w, err)
		return
	}

	write(w, subnet)
}

func (h *Handler) UpdateSubnet(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	subnet, err := h.ds.OneSubnet(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(subnet); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.UpdateSubnet(subnet); err != nil {
		writeError(w, err)
		return
	}

	write(w, subnet)
}

func (h *Handler) DeleteSubnet(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	if err := h.ds.DeleteSubnet(uuid); err != nil {
		writeError(w, err)
		return
	}

	writeDelete(w)
}
