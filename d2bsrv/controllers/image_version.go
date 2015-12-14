package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/imc-trading/dock2box/d2bsrv/models"
)

type ImageVersionController struct {
	database string
	session  *mgo.Session
}

func NewImageVersionController(s *mgo.Session) *ImageVersionController {
	return &ImageVersionController{
		database: "d2b",
		session:  s,
	}
}

func (c ImageVersionController) SetDatabase(database string) {
	c.database = database
}

func (c ImageVersionController) All(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	// Initialize empty struct
	s := models.Image{}

	// Get entry
	if err := c.session.DB(c.database).C("images").Find(bson.M{"image": name}).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s.Versions, http.StatusOK)
}

func (c ImageVersionController) AllByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get ID
	oid := bson.ObjectIdHex(id)

	// Initialize empty struct
	s := models.Image{}

	// Get entry
	if err := c.session.DB(c.database).C("images").FindId(oid).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s.Versions, http.StatusOK)
}
