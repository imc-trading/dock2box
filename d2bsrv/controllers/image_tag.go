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

type ImageTagController struct {
	database  string
	schemaURI string
	session   *mgo.Session
}

func NewImageTagController(s *mgo.Session) *ImageTagController {
	return &ImageTagController{
		database:  "d2b",
		schemaURI: "file://schemas/image-tag.json",
		session:   s,
	}
}

func (c *ImageTagController) SetDatabase(database string) {
	c.database = database
}

func (c *ImageTagController) SetSchemaURI(uri string) {
	c.schemaURI = uri + "/image-tag.json"
}

func (c *ImageTagController) CreateIndex() {
	index := mgo.Index{
		Key:    []string{"tag"},
		Unique: true,
	}

	if err := c.session.DB(c.database).C("image_tags").EnsureIndex(index); err != nil {
		panic(err)
	}
}

func (c *ImageTagController) All(w http.ResponseWriter, r *http.Request) {
	// Initialize empty struct list
	s := []models.ImageTag{}

	// Get all entries
	if err := c.session.DB(c.database).C("image_tags").Find(nil).All(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	/*
		if r.URL.Query().Get("embed") == "true" {
			for i, v := range s {
				// Get boot image
				if err := c.session.DB(c.database).C("boot_images").FindId(v.BootImageTagID).One(&s[i].BootImageTag); err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}
		}
	*/

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}

func (c *ImageTagController) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	// Initialize empty struct
	s := models.ImageTag{}

	// Get entry
	if err := c.session.DB(c.database).C("image_tags").Find(bson.M{"tag": name}).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	/*
		if r.URL.Query().Get("embed") == "true" {
			// Get boot image
			if err := c.session.DB(c.database).C("boot_images").FindId(s.BootImageTagID).One(&s.BootImageTag); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
	*/

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}

func (c *ImageTagController) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get ID
	oid := bson.ObjectIdHex(id)

	// Initialize empty struct
	s := models.ImageTag{}

	// Get entry
	if err := c.session.DB(c.database).C("image_tags").FindId(oid).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	/*
		if r.URL.Query().Get("embed") == "true" {
			// Get boot image
			if err := c.session.DB(c.database).C("boot_images").FindId(s.BootImageTagID).One(&s.BootImageTag); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
	*/

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}

func (c *ImageTagController) Create(w http.ResponseWriter, r *http.Request) {
	// Initialize empty struct
	s := models.ImageTag{}

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
	//	s.BootImageTagRef = "/boot-images/id/" + s.BootImageTagID.Hex()

	// Insert entry
	if err := c.session.DB(c.database).C("image_tags").Insert(s); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusCreated)
}

func (c *ImageTagController) Remove(w http.ResponseWriter, r *http.Request) {
	// Get name
	name := mux.Vars(r)["name"]

	// Remove entry
	if err := c.session.DB(c.database).C("image_tags").Remove(bson.M{"image_tags": name}); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	jsonWriter(w, r, nil, http.StatusOK)
}

func (c *ImageTagController) RemoveByID(w http.ResponseWriter, r *http.Request) {
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
	if err := c.session.DB(c.database).C("image_tags").RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	jsonWriter(w, r, nil, http.StatusOK)
}

func (c *ImageTagController) Update(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	// Initialize empty struct
	s := models.ImageTag{}

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
	//	s.BootImageTagRef = "/boot-images/id/" + s.BootImageTagID.Hex()

	// Update entry
	if err := c.session.DB(c.database).C("image_tags").Update(bson.M{"tag": name}, s); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK)
}
