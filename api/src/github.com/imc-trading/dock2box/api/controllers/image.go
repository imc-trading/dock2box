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

	"github.com/imc-trading/dock2box/api/models"
)

type ImageController struct {
	database  string
	schemaURI string
	session   *mgo.Session
	baseURI   string
	envelope  bool
	hateoas   bool
}

func NewImageController(s *mgo.Session, b string, e bool, h bool) *ImageController {
	return &ImageController{
		database:  "d2b",
		schemaURI: "file://schemas/image.json",
		session:   s,
		baseURI:   b,
		envelope:  e,
		hateoas:   h,
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
	es := reflect.ValueOf(models.Image{})
	keys, _ := structTags(es, "field", "bson")
	kinds, _ := structKinds(es, "field")

	// Query
	cond := bson.M{}
	qry := r.URL.Query()
	for k, v := range qry {
		if k == "envelope" || k == "embed" || k == "sort" || k == "hateoas" || k == "fields" {
			continue
		}
		//      fmt.Println(k, kinds[k]
		if _, ok := keys[k]; !ok {
			jsonError(w, r, fmt.Sprintf("Incorrect key used in query: %s", k), http.StatusBadRequest, c.envelope)
			return
		} else if kinds[k] == reflect.Bool {
			switch strings.ToLower(v[0]) {
			case "true":
				cond[keys[k]] = true
			case "false":
				cond[keys[k]] = false
			default:
				jsonError(w, r, fmt.Sprintf("Incorrect value for key used in query needs to be either true/false: %s", k), http.StatusBadRequest, c.envelope)
				return
			}
		} else if bson.IsObjectIdHex(v[0]) {
			cond[keys[k]] = bson.ObjectIdHex(v[0])
		} else {
			cond[keys[k]] = v[0]
		}
	}

	// Sort
	sort := []string{}
	if _, ok := qry["sort"]; ok {
		for _, str := range strings.Split(qry["sort"][0], ",") {
			op := ""
			k := str
			switch str[0] {
			case '+':
				op = "+"
				k = str[1:len(str)]
			case '-':
				op = "-"
				k = str[1:len(str)]
			}
			if _, ok := keys[k]; !ok {
				jsonError(w, r, fmt.Sprintf("Incorrect key used in sort: %s", k), http.StatusBadRequest, c.envelope)
				return
			}
			sort = append(sort, op+keys[k])
		}
	}

	// Fields
	fields := bson.M{}
	if _, ok := qry["fields"]; ok {
		for _, k := range strings.Split(qry["fields"][0], ",") {
			if _, ok := keys[k]; !ok {
				jsonError(w, r, fmt.Sprintf("Incorrect key used in fields: %s", k), http.StatusBadRequest, c.envelope)
				return
			}
			fields[keys[k]] = 1
		}
	}

	// Initialize empty struct list
	s := []models.Image{}

	// Get all entries
	if err := c.session.DB(c.database).C("images").Find(cond).Sort(sort...).Select(fields).All(&s); err != nil {
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

	// HATEOAS Links
	hateoas := c.hateoas
	switch strings.ToLower(r.URL.Query().Get("hateoas")) {
	case "true":
		hateoas = true
	case "false":
		hateoas = false
	}
	if hateoas == true {
		for i, v := range s {
			links := []models.Link{}

			if v.BootTagID != "" {
				links = append(links, models.Link{
					HRef:   c.baseURI + "/tags/" + v.BootTagID.Hex(),
					Rel:    "self",
					Method: "GET",
				})
			}

			s[i].Links = &links
		}
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK, c.envelope)
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

	// HATEOAS Links
	hateoas := c.hateoas
	switch strings.ToLower(r.URL.Query().Get("hateoas")) {
	case "true":
		hateoas = true
	case "false":
		hateoas = false
	}
	if hateoas == true {
		links := []models.Link{}

		if s.BootTagID != "" {
			links = append(links, models.Link{
				HRef:   c.baseURI + "/tags/" + s.BootTagID.Hex(),
				Rel:    "self",
				Method: "GET",
			})
		}

		s.Links = &links
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK, c.envelope)
}

func (c *ImageController) Create(w http.ResponseWriter, r *http.Request) {
	// Initialize empty struct
	s := models.Image{}

	// Decode JSON into struct
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		jsonError(w, r, "Failed to deconde JSON: "+err.Error(), http.StatusInternalServerError, c.envelope)
		return
	}

	// Set ID
	s.ID = bson.NewObjectId()

	// Validate input using JSON Schema
	docLoader := gojsonschema.NewGoLoader(s)
	schemaLoader := gojsonschema.NewReferenceLoader(c.schemaURI)

	res, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		jsonError(w, r, "Failed to load schema: "+err.Error(), http.StatusInternalServerError, c.envelope)
		return
	}

	if !res.Valid() {
		var errors []string
		for _, e := range res.Errors() {
			errors = append(errors, fmt.Sprintf("%s: %s", e.Context().String(), e.Description()))
		}
		jsonError(w, r, errors, http.StatusInternalServerError, c.envelope)
		return
	}

	// Insert entry
	if err := c.session.DB(c.database).C("images").Insert(s); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError, c.envelope)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusCreated, c.envelope)
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
		jsonError(w, r, "Failed to deconde JSON: "+err.Error(), http.StatusInternalServerError, c.envelope)
		return
	}

	// Validate input using JSON Schema
	docLoader := gojsonschema.NewGoLoader(s)
	schemaLoader := gojsonschema.NewReferenceLoader(c.schemaURI)

	res, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		jsonError(w, r, "Failed to load schema: "+err.Error(), http.StatusInternalServerError, c.envelope)
		return
	}

	if !res.Valid() {
		var errors []string
		for _, e := range res.Errors() {
			errors = append(errors, fmt.Sprintf("%s: %s", e.Context().String(), e.Description()))
		}
		jsonError(w, r, errors, http.StatusInternalServerError, c.envelope)
		return
	}

	// Update entry
	if err := c.session.DB(c.database).C("images").UpdateId(oid, s); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError, c.envelope)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK, c.envelope)
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

	// Initialize empty struct
	s := models.Image{}

	// Get entry
	if err := c.session.DB(c.database).C("images").FindId(oid).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Remove entry
	if err := c.session.DB(c.database).C("images").RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	jsonWriter(w, r, s, http.StatusOK, c.envelope)
}
