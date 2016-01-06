package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/imc-trading/dock2box/d2bsrv/models"
)

type ImageController struct {
	database  string
	schemaURI string
	session   *mgo.Session
}

func NewImageController(s *mgo.Session) *ImageController {
	return &ImageController{
		database:  "d2b",
		schemaURI: "file://schemas/image.json",
		session:   s,
	}
}

func (c *ImageController) SetDatabase(database string) {
	c.database = database
}

func (c *ImageController) SetSchemaURI(uri string) {
	c.schemaURI = uri + "/image.json"
}

func (c *ImageController) CreateIndex() {
	index := mgo.Index{
		Key:    []string{"image"},
		Unique: true,
	}

	if err := c.session.DB(c.database).C("images").EnsureIndex(index); err != nil {
		panic(err)
	}
}

func (c *ImageController) All(w http.ResponseWriter, r *http.Request) {
	// Get allowed key names
	keys, _ := structTags(reflect.ValueOf(models.Image{}), "json")

	// Query
	cond := bson.M{}
	qry := r.URL.Query()
	for k, v := range r.URL.Query() {
		if k == "envelope" || k == "embed" || k == "sort" {
			continue
		}
		if _, ok := keys[k]; !ok {
			jsonError(w, r, fmt.Sprintf("Incorrect key used in query: %s", k), http.StatusBadRequest)
			return
		} else if bson.IsObjectIdHex(v[0]) {
			cond[k] = bson.ObjectIdHex(v[0])
		} else {
			cond[k] = v[0]
		}
	}

	// Sort
	sort := []string{}
	if _, ok := qry["sort"]; ok {
		for _, k := range strings.Split(qry["sort"][0], ",") {
			if _, ok := keys[strings.TrimLeft(k, "+-")]; !ok {
				jsonError(w, r, fmt.Sprintf("Incorrect key used in sort: %s", k), http.StatusBadRequest)
				return
			}
			sort = append(sort, k)
		}
	}

	// Initialize empty struct list
	s := []models.Image{}

	// Get all entries
	if err := c.session.DB(c.database).C("images").Find(cond).Sort(sort...).All(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Embed related data
	if r.URL.Query().Get("embed") == "true" {
		for i, v := range s {
			if v.BootTagID != "" {
				// Get boot tag
				if err := c.session.DB(c.database).C("tags").FindId(v.BootTagID).One(&s[i].BootTag); err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}
		}
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}

func (c *ImageController) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get object id
	oid := bson.ObjectIdHex(id)

	// Initialize empty struct
	s := models.Image{}

	// Get entry
	if err := c.session.DB(c.database).C("images").FindId(oid).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Embed related data
	if r.URL.Query().Get("embed") == "true" {
		if s.BootTagID != "" {
			// Get boot tag
			if err := c.session.DB(c.database).C("tags").FindId(s.BootTagID).One(&s.BootTag); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}

func (c *ImageController) Create(w http.ResponseWriter, r *http.Request) {
	// Initialize empty struct
	s := models.Image{}

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

	// Insert entry
	if err := c.session.DB(c.database).C("images").Insert(s); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusCreated)
}

func (c *ImageController) Update(w http.ResponseWriter, r *http.Request) {
	// Get Id
	id := mux.Vars(r)["id"]

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get object id
	oid := bson.ObjectIdHex(id)

	// Initialize empty struct
	s := models.Image{}

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
	if err := c.session.DB(c.database).C("images").UpdateId(oid, s); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}

func (c *ImageController) Delete(w http.ResponseWriter, r *http.Request) {
	// Get Id
	id := mux.Vars(r)["id"]

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get object id
	oid := bson.ObjectIdHex(id)

	// Remove entry
	if err := c.session.DB(c.database).C("images").RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	jsonWriter(w, r, nil, http.StatusOK)
}
