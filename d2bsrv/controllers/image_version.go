package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/imc-trading/dock2box/d2bsrv/models"
)

type ImageVersionController struct {
	database  string
	schemaURI string
	session   *mgo.Session
}

func NewImageVersionController(s *mgo.Session) *ImageVersionController {
	return &ImageVersionController{
		database:  "d2b",
		schemaURI: "file://schemas/image.json",
		session:   s,
	}
}

func (c *ImageVersionController) SetDatabase(database string) {
	c.database = database
}

func (c *ImageVersionController) SetSchemaURI(uri string) {
	c.schemaURI = uri + "image.json"
}

func (c *ImageVersionController) All(w http.ResponseWriter, r *http.Request) {
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

func (c *ImageVersionController) AllByID(w http.ResponseWriter, r *http.Request) {
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

func (c *ImageVersionController) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	// Initialize empty struct
	s := models.Image{}

	// Get entry
	if err := c.session.DB(c.database).C("images").Find(bson.M{"image": name}).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, e := range s.Versions {
		if e.Version == version {
			// Write content-type, header and payload
			jsonWriter(w, r, e, http.StatusOK)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func (c *ImageVersionController) Create(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	// Initialize empty struct
	s := models.ImageVersion{}

	// Decode JSON into struct
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		jsonError(w, r, "Failed to deconde JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialize empty struct
	s2 := models.Image{}

	// Get entry
	if err := c.session.DB(c.database).C("images").Find(bson.M{"image": name}).One(&s2); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Validate input using JSON Schema
	docLoader := gojsonschema.NewGoLoader(s)
	schemaLoader := gojsonschema.NewReferenceLoader(c.schemaURI)

	res, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		jsonError(w, r, "Failed to load schema: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !res.Valid() {
		var errors []string
		for _, e := range res.Errors() {
			errors = append(errors, fmt.Sprintf("%s: %s", e.Context().String(), e.Description()))
		}
		jsonError(w, r, errors, http.StatusInternalServerError)
		return
	}

	for _, e := range s2.Versions {
		if e.Version == s.Version {
			jsonError(w, r, "Duplicate version already exists", http.StatusInternalServerError)
			return
		}
	}

	s2.Versions = append(s2.Versions, s)

	// Update entry
	if err := c.session.DB(c.database).C("images").Update(bson.M{"image": name}, s2); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusCreated)
}
