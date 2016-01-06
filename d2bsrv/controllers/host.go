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

type HostController struct {
	database  string
	schemaURI string
	session   *mgo.Session
}

func NewHostController(s *mgo.Session) *HostController {
	return &HostController{
		database:  "d2b",
		schemaURI: "file://schemas/host.json",
		session:   s,
	}
}

func (c *HostController) SetDatabase(database string) {
	c.database = database
}

func (c *HostController) SetSchemaURI(uri string) {
	c.schemaURI = uri + "/host.json"
}

func (c *HostController) CreateIndex() {
	index := mgo.Index{
		Key:    []string{"host"},
		Unique: true,
	}

	if err := c.session.DB(c.database).C("hosts").EnsureIndex(index); err != nil {
		panic(err)
	}
}

func (c *HostController) All(w http.ResponseWriter, r *http.Request) {
	// Initialize empty struct list
	s := []models.Host{}

	// Get all entries
	if err := c.session.DB(c.database).C("hosts").Find(nil).All(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}

func (c *HostController) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get object id
	oid := bson.ObjectIdHex(id)

	// Initialize empty struct
	s := models.Host{}

	// Get entry
	if err := c.session.DB(c.database).C("hosts").FindId(oid).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}

func (c *HostController) Create(w http.ResponseWriter, r *http.Request) {
	// Initialize empty struct
	s := models.Host{}

	// Decode JSON into struct
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		jsonError(w, r, "Failed to deconde JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set ID
	s.ID = bson.NewObjectId()

	// Validate input using JSON Schema
	docLoader := gojsonschema.NewGoLoader(s)
	schemaLoader := gojsonschema.NewReferenceLoader(c.schemaURI)

	res, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError)
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

	// Insert entry
	if err := c.session.DB(c.database).C("hosts").Insert(s); err != nil {
		jsonError(w, r, "Insert: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusCreated)
}

func (c *HostController) Update(w http.ResponseWriter, r *http.Request) {
	// Get ID
	id := mux.Vars(r)["id"]

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get object id
	oid := bson.ObjectIdHex(id)

	// Initialize empty struct
	s := models.Host{}

	// Decode JSON into struct
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		jsonError(w, r, "Failed to deconde JSON: "+err.Error(), http.StatusInternalServerError)
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

	// Update entry
	if err := c.session.DB(c.database).C("hosts").UpdateId(oid, s); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}

func (c *HostController) Delete(w http.ResponseWriter, r *http.Request) {
	// Get ID
	id := mux.Vars(r)["id"]

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get object id
	oid := bson.ObjectIdHex(id)

	// Initialize empty struct
	s := models.Host{}

	// Get entry
	if err := c.session.DB(c.database).C("hosts").FindId(oid).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Remove entry
	if err := c.session.DB(c.database).C("hosts").RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	jsonWriter(w, r, s, http.StatusOK)
}
