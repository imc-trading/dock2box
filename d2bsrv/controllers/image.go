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
	// Initialize empty struct list
	s := []models.Image{}

	// Get all entries
	if err := c.session.DB(c.database).C("images").Find(nil).All(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.URL.Query().Get("embed") == "true" {
		for i, v := range s {
			// Get boot image
			if err := c.session.DB(c.database).C("boot_images").FindId(v.BootImageID).One(&s[i].BootImage); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}

func (c *ImageController) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	// Initialize empty struct
	s := models.Image{}

	// Get entry
	if err := c.session.DB(c.database).C("images").Find(bson.M{"image": name}).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.URL.Query().Get("embed") == "true" {
		// Get boot image
		if err := c.session.DB(c.database).C("boot_images").FindId(s.BootImageID).One(&s.BootImage); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}

func (c *ImageController) GetByID(w http.ResponseWriter, r *http.Request) {
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

	if r.URL.Query().Get("embed") == "true" {
		// Get boot image
		if err := c.session.DB(c.database).C("boot_images").FindId(s.BootImageID).One(&s.BootImage); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
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

	// Set refs
	s.BootImageRef = "/boot-images/id/" + s.BootImageID.Hex()

	// Insert entry
	if err := c.session.DB(c.database).C("images").Insert(s); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusCreated)
}

func (c *ImageController) Remove(w http.ResponseWriter, r *http.Request) {
	// Get name
	name := mux.Vars(r)["name"]

	// Remove entry
	if err := c.session.DB(c.database).C("images").Remove(bson.M{"image": name}); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	jsonWriter(w, r, nil, http.StatusOK)
}

func (c *ImageController) RemoveByID(w http.ResponseWriter, r *http.Request) {
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
	if err := c.session.DB(c.database).C("images").RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	jsonWriter(w, r, nil, http.StatusOK)
}

func (c *ImageController) Update(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

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

	// Set refs
	s.BootImageRef = "/boot-images/id/" + s.BootImageID.Hex()

	// Update entry
	if err := c.session.DB(c.database).C("images").Update(bson.M{"image": name}, s); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}
