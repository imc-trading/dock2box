package handler

import (
	"encoding/json"
	"net/http"

	"github.com/imc-trading/dock2box/model"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	l := &model.Login{}
	if err := json.NewDecoder(r.Body).Decode(l); err != nil {
		writeError(w, err)
		return
	}
	defer r.Body.Close()

	u, err := h.conn.Login(l.Username, l.Password)
	if err != nil {
		writeError(w, err)
		return
	}

	s, err := h.jwt.NewToken(u).Sign()
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(s))
}

func (h *Handler) Renew(w http.ResponseWriter, r *http.Request) {
	t, err := h.jwt.ParseTokenHeader(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(err.Error()))
		return
	}

	s, err := t.Renew().Sign()
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(s))
}

func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	t, err := h.jwt.ParseTokenHeader(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(err.Error()))
		return
	}

	write(w, t.Claims)
}
