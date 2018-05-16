package handler

import (
	"net/http"

	"github.com/mickep76/encdec"
	_ "github.com/mickep76/encdec/json"
	"github.com/mickep76/qry"
)

func (h *Handler) AllServers(w http.ResponseWriter, r *http.Request) {
	all, err := h.ds.AllServers()
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var b []byte
	vals := r.URL.Query()
	if len(vals) > 0 {
		q, err := qry.FromURL(r.URL.Query())
		if err != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		filtered, err := q.Query(all)
		if err != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		b, _ = encdec.ToBytes("json", filtered, encdec.WithIndent("  "))
	} else {
		b, _ = encdec.ToBytes("json", all, encdec.WithIndent("  "))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
