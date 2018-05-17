package handler

import (
	"net/http"

	"github.com/mickep76/qry"
)

func (h *Handler) AllServers(w http.ResponseWriter, r *http.Request) {
	q, err := qry.FromURL(r.URL.Query())
	if err != nil {
		writeError(w, err)
		return
	}

	servers, err := h.ds.QueryServers(q)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, servers)
}
