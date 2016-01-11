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

type SubnetController struct {
	database  string
	schemaURI string
	session   *mgo.Session
	baseURI   string
	envelope  bool
	hateoas   bool
}

func NewSubnetController(s *mgo.Session, b string, e bool, h bool) *SubnetController {
	return &SubnetController{
		database:  "d2b",
		schemaURI: "file://schemas/subnet.json",
		session:   s,
		baseURI:   b,
		envelope:  e,
		hateoas:   h,
	}
}

func (c *SubnetController) SetDatabase(database string) {
	c.database = database
}

func (c *SubnetController) SetSchemaURI(uri string) {
	c.schemaURI = uri + "/subnet.json"
}

func (c *SubnetController) CreateIndex() {
	index := mgo.Index{
		Key:    []string{"subnet"},
		Unique: true,
	}

	if err := c.session.DB(c.database).C("subnets").EnsureIndex(index); err != nil {
		panic(err)
	}
}

func (c *SubnetController) All(w http.ResponseWriter, r *http.Request) {
	// Get allowed key names
	keys, _ := structTags(reflect.ValueOf(models.Tag{}), "field", "bson")

	// Query
	cond := bson.M{}
	qry := r.URL.Query()
	for k, v := range qry {
		if k == "envelope" || k == "embed" || k == "sort" || k == "hateoas" || k == "fields" {
			continue
		}
		if _, ok := keys[k]; !ok {
			jsonError(w, r, fmt.Sprintf("Incorrect key used in query: %s", k), http.StatusBadRequest, c.envelope)
			return
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
	s := []models.Subnet{}

	// Get all entries
	if err := c.session.DB(c.database).C("subnets").Find(cond).Sort(sort...).Select(fields).All(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Embed related data
	if r.URL.Query().Get("embed") == "true" {
		for i, v := range s {
			// Get site
			if err := c.session.DB(c.database).C("sites").FindId(v.SiteID).One(&s[i].Site); err != nil {
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
		for i, v := range s {
			s[i].Links = &[]models.Link{
				models.Link{
					HRef:   c.baseURI + "/sites/" + v.SiteID.Hex(),
					Rel:    "self",
					Method: "GET",
				},
			}
		}
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK, c.envelope)
}

func (c *SubnetController) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get object id
	oid := bson.ObjectIdHex(id)

	// Initialize empty struct
	s := models.Subnet{}

	// Get entry
	if err := c.session.DB(c.database).C("subnets").FindId(oid).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.URL.Query().Get("embed") == "true" {
		// Get site
		if err := c.session.DB(c.database).C("sites").FindId(s.SiteID).One(&s.Site); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	// Embed related data
	if r.URL.Query().Get("embed") == "true" {
		// Get site
		if err := c.session.DB(c.database).C("sites").FindId(s.SiteID).One(&s.Site); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
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
		s.Links = &[]models.Link{
			models.Link{
				HRef:   c.baseURI + "/sites/" + s.SiteID.Hex(),
				Rel:    "self",
				Method: "GET",
			},
		}
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK, c.envelope)
}

func (c *SubnetController) Create(w http.ResponseWriter, r *http.Request) {
	// Initialize empty struct
	s := models.Subnet{}

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
		jsonError(w, r, err.Error(), http.StatusInternalServerError, c.envelope)
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
	if err := c.session.DB(c.database).C("subnets").Insert(s); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError, c.envelope)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusCreated, c.envelope)
}

func (c *SubnetController) Update(w http.ResponseWriter, r *http.Request) {
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
	s := models.Subnet{}

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
	if err := c.session.DB(c.database).C("subnets").UpdateId(oid, s); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError, c.envelope)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK, c.envelope)
}

func (c *SubnetController) Delete(w http.ResponseWriter, r *http.Request) {
	// Get ID
	id := mux.Vars(r)["id"]

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get object id
	oid := bson.ObjectIdHex(id)

	// Remove entry
	if err := c.session.DB(c.database).C("subnets").RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	jsonWriter(w, r, nil, http.StatusOK, c.envelope)
}
