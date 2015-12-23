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
		schemaURI: "file://schemas/image-version.json",
		session:   s,
	}
}

func (c *ImageVersionController) SetDatabase(database string) {
	c.database = database
}

func (c *ImageVersionController) SetSchemaURI(uri string) {
	c.schemaURI = uri + "/image-version.json"
}

func (c *ImageVersionController) CreateIndex() {
	index := mgo.Index{
		Key:    []string{"version"},
		Unique: true,
	}

	if err := c.session.DB(c.database).C("image_versions").EnsureIndex(index); err != nil {
		panic(err)
	}
}

func (c *ImageVersionController) All(w http.ResponseWriter, r *http.Request) {
	// Initialize empty struct list
	s := []models.ImageVersion{}

	// Get all entries
	if err := c.session.DB(c.database).C("image_versions").Find(nil).All(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	/*
		if r.URL.Query().Get("embed") == "true" {
			for i, v := range s {
				// Get boot image
				if err := c.session.DB(c.database).C("boot_images").FindId(v.BootImageVersionID).One(&s[i].BootImageVersion); err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}
		}
	*/

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}

func (c *ImageVersionController) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	// Initialize empty struct
	s := models.ImageVersion{}

	// Get entry
	if err := c.session.DB(c.database).C("image_versions").Find(bson.M{"version": name}).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	/*
		if r.URL.Query().Get("embed") == "true" {
			// Get boot image
			if err := c.session.DB(c.database).C("boot_images").FindId(s.BootImageVersionID).One(&s.BootImageVersion); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
	*/

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}

func (c *ImageVersionController) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get ID
	oid := bson.ObjectIdHex(id)

	// Initialize empty struct
	s := models.ImageVersion{}

	// Get entry
	if err := c.session.DB(c.database).C("image_versions").FindId(oid).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	/*
		if r.URL.Query().Get("embed") == "true" {
			// Get boot image
			if err := c.session.DB(c.database).C("boot_images").FindId(s.BootImageVersionID).One(&s.BootImageVersion); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
	*/

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}

func (c *ImageVersionController) Create(w http.ResponseWriter, r *http.Request) {
	// Initialize empty struct
	s := models.ImageVersion{}

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
	//	s.BootImageVersionRef = "/boot-images/id/" + s.BootImageVersionID.Hex()

	// Insert entry
	if err := c.session.DB(c.database).C("image_versions").Insert(s); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusCreated)
}

func (c *ImageVersionController) Remove(w http.ResponseWriter, r *http.Request) {
	// Get name
	name := mux.Vars(r)["name"]

	// Remove entry
	if err := c.session.DB(c.database).C("image_versions").Remove(bson.M{"image_versions": name}); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	jsonWriter(w, r, nil, http.StatusOK)
}

func (c *ImageVersionController) RemoveByID(w http.ResponseWriter, r *http.Request) {
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
	if err := c.session.DB(c.database).C("image_versions").RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	jsonWriter(w, r, nil, http.StatusOK)
}

func (c *ImageVersionController) Update(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	// Initialize empty struct
	s := models.ImageVersion{}

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
	//	s.BootImageVersionRef = "/boot-images/id/" + s.BootImageVersionID.Hex()

	// Update entry
	if err := c.session.DB(c.database).C("image_versions").Update(bson.M{"version": name}, s); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}
