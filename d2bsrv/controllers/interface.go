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

type InterfaceController struct {
	database  string
	schemaURI string
	session   *mgo.Session
	baseURI   string
	envelope  bool
	hateoas   bool
}

func NewInterfaceController(s *mgo.Session, b string, e bool, h bool) *InterfaceController {
	return &InterfaceController{
		database:  "d2b",
		schemaURI: "file://schemas/interface.json",
		session:   s,
		baseURI:   b,
		envelope:  e,
		hateoas:   h,
	}
}

func (c *InterfaceController) SetDatabase(database string) {
	c.database = database
}

func (c *InterfaceController) SetSchemaURI(uri string) {
	c.schemaURI = uri + "/interface.json"
}

func (c *InterfaceController) CreateIndex() {
	index := mgo.Index{
		Key:    []string{"hostId", "interface"},
		Unique: true,
	}

	if err := c.session.DB(c.database).C("interfaces").EnsureIndex(index); err != nil {
		panic(err)
	}

	index = mgo.Index{
		Key:    []string{"hwAddr"},
		Unique: true,
	}

	if err := c.session.DB(c.database).C("interfaces").EnsureIndex(index); err != nil {
		panic(err)
	}
}

func (c *InterfaceController) All(w http.ResponseWriter, r *http.Request) {
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
	s := []models.Interface{}

	// Get all entries
	if err := c.session.DB(c.database).C("interfaces").Find(cond).Sort(sort...).Select(fields).All(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Embed related data
	if r.URL.Query().Get("embed") == "true" {
		for i, v := range s {
			if v.SubnetID != "" {
				// Get subnet
				if err := c.session.DB(c.database).C("subnets").FindId(v.SubnetID).One(&s[i].Subnet); err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}
			// Get Host
			if err := c.session.DB(c.database).C("hosts").FindId(v.HostID).One(&s[i].Host); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK, c.envelope)
}

func (c *InterfaceController) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Validate ObjectId
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get object id
	oid := bson.ObjectIdHex(id)

	// Initialize empty struct
	s := models.Interface{}

	// Get entry
	if err := c.session.DB(c.database).C("interfaces").FindId(oid).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Embed related data
	if r.URL.Query().Get("embed") == "true" {
		if s.SubnetID != "" {
			// Get subnet
			if err := c.session.DB(c.database).C("subnets").FindId(s.SubnetID).One(&s.Subnet); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
		// Get Host
		if err := c.session.DB(c.database).C("hosts").FindId(s.HostID).One(&s.Host); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK, c.envelope)
}

func (c *InterfaceController) Create(w http.ResponseWriter, r *http.Request) {
	// Initialize empty struct
	s := models.Interface{}

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
	if err := c.session.DB(c.database).C("interfaces").Insert(s); err != nil {
		jsonError(w, r, "Insert: "+err.Error(), http.StatusInternalServerError, c.envelope)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusCreated, c.envelope)
}

func (c *InterfaceController) Update(w http.ResponseWriter, r *http.Request) {
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
	s := models.Interface{}

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
	if err := c.session.DB(c.database).C("interfaces").UpdateId(oid, s); err != nil {
		jsonError(w, r, err.Error(), http.StatusInternalServerError, c.envelope)
		return
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK, c.envelope)
}

func (c *InterfaceController) Delete(w http.ResponseWriter, r *http.Request) {
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
	s := models.Interface{}

	// Get entry
	if err := c.session.DB(c.database).C("interfaces").FindId(oid).One(&s); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Remove entry
	if err := c.session.DB(c.database).C("interfaces").RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	jsonWriter(w, r, s, http.StatusOK, c.envelope)
}
