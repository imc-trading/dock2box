package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mickep76/qry"

	"github.com/imc-trading/dock2box/model"
)

func (h *Handler) AllRoles(w http.ResponseWriter, r *http.Request) {
	q, err := qry.FromURL(r.URL.Query())
	if err != nil {
		writeError(w, err)
		return
	}

	roles, err := h.ds.QueryRoles(q)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, roles)
}

func (h *Handler) OneRole(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	role, err := h.ds.OneHost(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, role)
}

func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	role := &model.Role{}
	if err := json.NewDecoder(r.Body).Decode(role); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.CreateRole(role); err != nil {
		writeError(w, err)
		return
	}

	write(w, role)
}

func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	role, err := h.ds.OneRole(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(role); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.UpdateRole(role); err != nil {
		writeError(w, err)
		return
	}

	write(w, role)
}

func (h *Handler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	if err := h.ds.DeleteRole(uuid); err != nil {
		writeError(w, err)
		return
	}

	writeDelete(w)
}
