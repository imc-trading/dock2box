package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mickep76/qry"

	"github.com/imc-trading/dock2box/model"
)

func (h *Handler) AllImages(w http.ResponseWriter, r *http.Request) {
	q, err := qry.FromURL(r.URL.Query())
	if err != nil {
		writeError(w, err)
		return
	}

	images, err := h.ds.QueryImages(q)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, images)
}

func (h *Handler) OneImage(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	image, err := h.ds.OneHost(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, image)
}

func (h *Handler) CreateImage(w http.ResponseWriter, r *http.Request) {
	image := &model.Image{}
	if err := json.NewDecoder(r.Body).Decode(image); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.CreateImage(image); err != nil {
		writeError(w, err)
		return
	}

	write(w, image)
}

func (h *Handler) UpdateImage(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	image, err := h.ds.OneImage(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(image); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.UpdateImage(image); err != nil {
		writeError(w, err)
		return
	}

	write(w, image)
}

func (h *Handler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	if err := h.ds.DeleteImage(uuid); err != nil {
		writeError(w, err)
		return
	}

	writeDelete(w)
}
