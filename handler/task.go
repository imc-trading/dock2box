package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mickep76/qry"

	"github.com/imc-trading/dock2box/model"
)

func (h *Handler) AllTasks(w http.ResponseWriter, r *http.Request) {
	q, err := qry.FromURL(r.URL.Query())
	if err != nil {
		writeError(w, err)
		return
	}

	tasks, err := h.ds.QueryTasks(q)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, tasks)
}

func (h *Handler) OneTask(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	task, err := h.ds.OneHost(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	write(w, task)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	task := &model.Task{}
	if err := json.NewDecoder(r.Body).Decode(task); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.CreateTask(task); err != nil {
		writeError(w, err)
		return
	}

	write(w, task)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	task, err := h.ds.OneTask(uuid)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(task); err != nil {
		writeError(w, err)
		return
	}

	if err := h.ds.UpdateTask(task); err != nil {
		writeError(w, err)
		return
	}

	write(w, task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]

	if err := h.ds.DeleteTask(uuid); err != nil {
		writeError(w, err)
		return
	}

	writeDelete(w)
}
