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

type SiteController struct {
	database  string
	schemaURI string
	session   *mgo.Session
}

func NewSiteController(s *mgo.Session) *SiteController {
	return &SiteController{
		database:  "d2b",
		schemaURI: "file://schemas/site.json",
		session:   s,
	}
}

func (c *SiteController) SetDatabase(database string) {
	c.database = database
}

func (c *SiteController) SetSchemaURI(uri string) {
	c.schemaURI = uri + "/site.json"
}

func (c *SiteController) CreateIndex() {
	index := mgo.Index{
		Key:    []string{"site"},
		Unique: true,
	}

	if err := c.session.DB(c.database).C("sites").EnsureIndex(index); err != nil {
		panic(err)
	}
}

func (c *SiteController) All(w http.ResponseWriter, r *http.Request) {
	// Initialize empty struct list
	s := []models.Site{}

	// Get all entries
	if err := c.session.DB(c.database).C("sites").Find(nil).All(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}

func (c *SiteController) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	// Initialize empty struct
	s := models.Site{}

	// Get entry
	if err := c.session.DB(c.database).C("sites").Find(bson.M{"site": name}).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}

func (c *SiteController) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get ID
	oid := bson.ObjectIdHex(id)

	// Initialize empty struct
	s := models.Site{}

	// Get entry
	if err := c.session.DB(c.database).C("sites").FindId(oid).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}

func (c *SiteController) Create(w http.ResponseWriter, r *http.Request) {
	// Initialize empty struct
	s := models.Site{}

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
	if err := c.session.DB(c.database).C("sites").Insert(s); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusCreated)
}

func (c *SiteController) Remove(w http.ResponseWriter, r *http.Request) {
	// Get name
	name := mux.Vars(r)["name"]

	// Remove entry
	if err := c.session.DB(c.database).C("sites").Remove(bson.M{"site": name}); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	jsonWriter(w, r, nil, http.StatusOK)
}

func (c *SiteController) RemoveByID(w http.ResponseWriter, r *http.Request) {
	// Get ID
	id := mux.Vars(r)["id"]

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get new ID
	oid := bson.ObjectIdHex(id)

	// Remove entry
	if err := c.session.DB(c.database).C("sites").RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	jsonWriter(w, r, nil, http.StatusOK)
}

func (c *SiteController) Update(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	// Initialize empty struct
	s := models.Site{}

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
	if err := c.session.DB(c.database).C("sites").Update(bson.M{"site": name}, s); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}
