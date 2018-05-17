package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mickep76/qry"

	"github.com/imc-trading/dock2box/model"
)

func (h *Handler) AllTenants(w http.ResponseWriter, r *http.Request) {
	q, err := qry.FromURL(r.URL.Query())
	if err != nil {
		writeError(w, err)
		return
	}

	tenants, err := h.ds.QueryTenants(q)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, tenants)
}

func (h *Handler) OneTenant(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	tenant, err := h.ds.OneHost(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, tenant)
}

func (h *Handler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	tenant := &model.Tenant{}
	if err := json.NewDecoder(r.Body).Decode(tenant); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.CreateTenant(tenant); err != nil {
		writeError(w, err)
		return
	}

	write(w, tenant)
}

func (h *Handler) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	tenant, err := h.ds.OneTenant(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(tenant); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.UpdateTenant(tenant); err != nil {
		writeError(w, err)
		return
	}

	write(w, tenant)
}
