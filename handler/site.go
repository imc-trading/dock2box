package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mickep76/qry"

	"github.com/imc-trading/dock2box/model"
)

func (h *Handler) AllSites(w http.ResponseWriter, r *http.Request) {
	q, err := qry.FromURL(r.URL.Query())
	if err != nil {
		writeError(w, err)
		return
	}

	sites, err := h.ds.QuerySites(q)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, sites)
}

func (h *Handler) OneSite(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	site, err := h.ds.OneHost(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, site)
}

func (h *Handler) CreateSite(w http.ResponseWriter, r *http.Request) {
	site := &model.Site{}
	if err := json.NewDecoder(r.Body).Decode(site); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.CreateSite(site); err != nil {
		writeError(w, err)
		return
	}

	write(w, site)
}

func (h *Handler) UpdateSite(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	site, err := h.ds.OneSite(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(site); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.UpdateSite(site); err != nil {
		writeError(w, err)
		return
	}

	write(w, site)
}
