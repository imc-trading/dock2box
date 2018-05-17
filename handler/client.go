package handler

import (
	"net/http"

	"github.com/mickep76/qry"
)

func (h *Handler) AllClients(w http.ResponseWriter, r *http.Request) {
	q, err := qry.FromURL(r.URL.Query())
	if err != nil {
		writeError(w, err)
		return
	}

	clients, err := h.ds.QueryClients(q)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, clients)
}
