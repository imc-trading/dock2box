package controllers

import (
	"encoding/json"
	"net/http"
)

type JSONEnvelope struct {
	Code   int         `json:"code"`
	Data   interface{} `json:"data,omitempty"`
	Errors interface{} `json:"errors,omitempty"`
}

func jsonError(w http.ResponseWriter, r *http.Request, e interface{}, c int, envelope bool) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(c)

	if envelope == false || r.URL.Query().Get("envelope") == "false" {
		var b []byte
		if r.URL.Query().Get("indent") == "false" {
			b, _ = json.Marshal(&e)
		} else {
			b, _ = json.MarshalIndent(&e, "", "  ")
		}
		w.Write(b)
		return
	}

	s := JSONEnvelope{
		Code:   c,
		Errors: e,
	}

	var b []byte
	if r.URL.Query().Get("indent") == "false" {
		b, _ = json.Marshal(&s)
	} else {
		b, _ = json.MarshalIndent(&s, "", "  ")
	}
	w.Write(b)
}

func jsonWriter(w http.ResponseWriter, r *http.Request, d interface{}, c int, envelope bool) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(c)

	if envelope == false || r.URL.Query().Get("envelope") == "false" {
		var b []byte
		if r.URL.Query().Get("indent") == "false" {
			b, _ = json.Marshal(&d)
		} else {
			b, _ = json.MarshalIndent(&d, "", "  ")
		}
		w.Write(b)
		return
	}

	s := JSONEnvelope{
		Code: c,
		Data: d,
	}

	var b []byte
	if r.URL.Query().Get("indent") == "false" {
		b, _ = json.Marshal(&s)
	} else {
		b, _ = json.MarshalIndent(&s, "", "  ")
	}
	w.Write(b)
}
