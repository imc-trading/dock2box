package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"

	"github.com/imc-trading/dock2box/d2bsrv/controllers"
	"github.com/imc-trading/dock2box/d2bsrv/version"
)

func main() {
	// Options.
	appVersion := flag.Bool("version", false, "Version")
	bind := flag.String("bind", "127.0.0.1:8080", "Bind to address and port")
	database := flag.String("database", "d2b", "Database name")
	schemaURI := flag.String("schema-uri", "file://schemas", "URI to JSON schemas")
	flag.Parse()

	// Print version.
	if *appVersion {
		fmt.Printf("d2bsrv %s\n", version.Version)
		os.Exit(0)
	}

	log.Printf("Using JSON schema URI: %s", *schemaURI)

	// Create new router
	r := mux.NewRouter()

	// Host
	// Get Controller instance
	host := controllers.NewHostController(getSession())

	// Set Database
	host.SetDatabase(*database)

	// Set Schema URI
	host.SetSchemaURI(*schemaURI)

	// Create Index
	host.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/hosts", host.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/hosts/{name}", host.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/hosts/id/{id}", host.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/hosts", host.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/hosts/{name}", host.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/hosts/{name}", host.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/hosts/id/{id}", host.RemoveByID).Methods("DELETE")

	// Interface
	// Get Controller instance
	intfs := controllers.NewInterfaceController(getSession())

	// Set Database
	intfs.SetDatabase(*database)

	// Set Schema URI
	intfs.SetSchemaURI(*schemaURI)

	// Create Index
	//    host.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/interfaces", intfs.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/interfaces/{name}", intfs.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/interfaces/id/{id}", intfs.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/interfaces", intfs.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/interfaces/{name}", intfs.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/interfaces/{name}", intfs.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/interfaces/id/{id}", intfs.RemoveByID).Methods("DELETE")

	// Site
	// Get Controller instance
	site := controllers.NewSiteController(getSession())

	// Set Database
	site.SetDatabase(*database)

	// Set Schema URI
	site.SetSchemaURI(*schemaURI)

	// Create Index
	site.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/sites", site.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/sites/{name}", site.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/sites/id/{id}", site.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/sites", site.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/sites/{name}", site.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/sites/{name}", site.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/sites/id/{id}", site.RemoveByID).Methods("DELETE")

	// Subnet
	// Get Controller instance
	subnet := controllers.NewSubnetController(getSession())

	// Set Schema URI
	subnet.SetSchemaURI(*schemaURI)

	// Set Database
	subnet.SetDatabase(*database)

	// Create Index
	subnet.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/subnets", subnet.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/subnets/{name}-{prefix}", subnet.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/subnets/id/{id}", subnet.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/subnets", subnet.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/subnets/{name}-{prefix}", subnet.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/subnets/{name}-{prefix}", subnet.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/subnets/id/{id}", subnet.RemoveByID).Methods("DELETE")

	// Image
	// Get Controller instance
	image := controllers.NewImageController(getSession())

	// Set Database
	image.SetDatabase(*database)

	// Set Schema URI
	image.SetSchemaURI(*schemaURI)

	// Create Index
	image.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/images", image.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/images/{name}", image.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/images/id/{id}", image.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/images", image.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/images/{name}", image.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/images/{name}", image.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/images/id/{id}", image.RemoveByID).Methods("DELETE")

	// Tag
	// Get Controller instance
	tag := controllers.NewTagController(getSession())

	// Set Database
	tag.SetDatabase(*database)

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/tags", tag.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/tags/{name}", tag.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/tags/id/{id}", tag.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/tags", tag.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/tags/{name}", tag.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/tags/{name}", tag.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/tags/id/{id}", tag.RemoveByID).Methods("DELETE")

	// Tenant
	// Get Controller instance
	tenant := controllers.NewTenantController(getSession())

	// Set Database
	tenant.SetDatabase(*database)

	// Set Schema URI
	tenant.SetSchemaURI(*schemaURI)

	// Create Index
	tenant.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/tenants", tenant.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/tenants/{name}", tenant.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/tenants/id/{id}", tenant.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/tenants", tenant.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/tenants/{name}", tenant.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/tenants/{name}", tenant.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/tenants/id/{id}", tenant.RemoveByID).Methods("DELETE")

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
