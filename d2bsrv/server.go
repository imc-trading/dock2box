package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"

	"github.com/imc-trading/dock2box/d2bsrv/controllers"
	"github.com/imc-trading/dock2box/d2bsrv/version"
)

func main() {
	// Options
	appVersion := flag.Bool("version", false, "Version")
	bind := flag.String("bind", "0.0.0.0:8080", "Bind to address and port")
	database := flag.String("database", "d2b", "Database name")
	schemaURI := flag.String("schema-uri", "file://schemas", "URI to JSON schemas")
	baseURI := flag.String("base-uri", "", "Base URI for server")
	disableHATEOAS := flag.Bool("disable-hateoas", false, "Disable HATEOAS per default")
	enableEnvelope := flag.Bool("enable-envelope", false, "Enable Envelopes per default")
	flag.Parse()

	// Print version
	if *appVersion {
		fmt.Printf("d2bsrv %s\n", version.Version)
		os.Exit(0)
	}

	log.Printf("Using JSON schema URI: %s", *schemaURI)

	// Get Base URI
	if *baseURI == "" {
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatal(err.Error())
		}
		port := strings.Split(*bind, ":")[1]
		str := "http://" + hostname + ":" + port + "/" + version.APIVersion
		baseURI = &str
	}

	log.Printf("Using Base URI: %s", *baseURI)

	// Get global HATEOAS setting
	hateoas := true
	if *disableHATEOAS == true {
		hateoas = false
	}
	log.Printf("Use HATEOAS per default: %v", hateoas)

	// Get global envelope setting
	envelope := false
	if *enableEnvelope == true {
		envelope = true
	}
	log.Printf("Use Envelope per default: %v", envelope)

	// Create new router
	r := mux.NewRouter()

	// Root
	// Get Controller instance
	root := controllers.NewRootController(*baseURI, envelope, hateoas)

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion, root.All).Methods("GET")

	// Host
	// Get Controller instance
	host := controllers.NewHostController(getSession(), *baseURI, envelope, hateoas)

	// Set Database
	host.SetDatabase(*database)

	// Set Schema URI
	host.SetSchemaURI(*schemaURI)

	// Create Index
	host.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/hosts", host.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/hosts/{id}", host.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/hosts", host.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/hosts/{id}", host.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/hosts/{id}", host.Delete).Methods("DELETE")

	// Interface
	// Get Controller instance
	intfs := controllers.NewInterfaceController(getSession(), *baseURI, envelope, hateoas)

	// Set Database
	intfs.SetDatabase(*database)

	// Set Schema URI
	intfs.SetSchemaURI(*schemaURI)

	// Create Index
	intfs.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/interfaces", intfs.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/interfaces/{id}", intfs.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/interfaces", intfs.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/interfaces/{id}", intfs.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/interfaces/{id}", intfs.Delete).Methods("DELETE")

	// Site
	// Get Controller instance
	site := controllers.NewSiteController(getSession(), *baseURI, envelope, hateoas)

	// Set Database
	site.SetDatabase(*database)

	// Set Schema URI
	site.SetSchemaURI(*schemaURI)

	// Create Index
	site.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/sites", site.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/sites/{id}", site.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/sites", site.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/sites/{id}", site.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/sites/{id}", site.Delete).Methods("DELETE")

	// Subnet
	// Get Controller instance
	subnet := controllers.NewSubnetController(getSession(), *baseURI, envelope, hateoas)

	// Set Schema URI
	subnet.SetSchemaURI(*schemaURI)

	// Set Database
	subnet.SetDatabase(*database)

	// Create Index
	subnet.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/subnets", subnet.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/subnets/{id}", subnet.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/subnets", subnet.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/subnets/{id}", subnet.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/subnets/{id}", subnet.Delete).Methods("DELETE")

	// Image
	// Get Controller instance
	image := controllers.NewImageController(getSession(), *baseURI, envelope, hateoas)

	// Set Database
	image.SetDatabase(*database)

	// Set Schema URI
	image.SetSchemaURI(*schemaURI)

	// Create Index
	image.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/images", image.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/images/{id}", image.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/images", image.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/images/{id}", image.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/images/{id}", image.Delete).Methods("DELETE")

	// Tag
	// Get Controller instance
	tag := controllers.NewTagController(getSession(), *baseURI, envelope, hateoas)

	// Set Database
	tag.SetDatabase(*database)

	// Create Index
	tag.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/tags", tag.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/tags/{id}", tag.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/tags", tag.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/tags/{id}", tag.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/tags/{id}", tag.Delete).Methods("DELETE")

	// Tenant
	// Get Controller instance
	tenant := controllers.NewTenantController(getSession(), *baseURI, envelope, hateoas)

	// Set Database
	tenant.SetDatabase(*database)

	// Set Schema URI
	tenant.SetSchemaURI(*schemaURI)

	// Create Index
	tenant.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/tenants", tenant.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/tenants/{id}", tenant.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/tenants", tenant.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/tenants/{id}", tenant.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/tenants/{name}", tenant.Delete).Methods("DELETE")

	// PXE Menu
	// Get Controller instance
	pc := controllers.NewPXEMenuController(getSession())

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/ipxe/{hwaddr}", pc.PXEMenu).Methods("GET")

	// Static Files
	// Schemas
	schemas := http.StripPrefix("/"+version.APIVersion+"/schemas/", http.FileServer(http.Dir("schemas")))
	r.PathPrefix("/" + version.APIVersion + "/schemas/").Handler(schemas)

	// Images
	img := http.StripPrefix("/img/", http.FileServer(http.Dir("img")))
	r.PathPrefix("/img/").Handler(img)

	http.Handle("/", r)

	// Fire up the server
	log.Printf("Bind to: %s", *bind)
	logr := handlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(*bind, logr))
}

// getSession creates a new mongo session and panics if connection error occurs
func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}

	// Deliver session
	return s
}
