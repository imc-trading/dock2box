package handler

import (
	"net/http"

	"github.com/mickep76/encdec"
	_ "github.com/mickep76/encdec/json"
	"github.com/mickep76/qry"
)

func (h *Handler) AllClients(w http.ResponseWriter, r *http.Request) {
	q, err := qry.FromURL(r.URL.Query())
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	v, err := h.ds.QueryClients(q)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	b, _ := encdec.ToBytes("json", v, encdec.WithIndent("  "))
	w.Write(b)
}
