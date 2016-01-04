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
	hc := controllers.NewHostController(getSession())

	// Set Database
	hc.SetDatabase(*database)

	// Set Schema URI
	hc.SetSchemaURI(*schemaURI)

	// Create Index
	hc.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/hosts", hc.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/hosts/{name}", hc.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/hosts/id/{id}", hc.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/hosts", hc.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/hosts/{name}", hc.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/hosts/{name}", hc.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/hosts/id/{id}", hc.RemoveByID).Methods("DELETE")

	// Site
	// Get Controller instance
	sc := controllers.NewSiteController(getSession())

	// Set Database
	sc.SetDatabase(*database)

	// Set Schema URI
	sc.SetSchemaURI(*schemaURI)

	// Create Index
	sc.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/sites", sc.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/sites/{name}", sc.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/sites/id/{id}", sc.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/sites", sc.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/sites/{name}", sc.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/sites/{name}", sc.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/sites/id/{id}", sc.RemoveByID).Methods("DELETE")

	// Subnet
	// Get Controller instance
	suc := controllers.NewSubnetController(getSession())

	// Set Schema URI
	suc.SetSchemaURI(*schemaURI)

	// Set Database
	suc.SetDatabase(*database)

	// Create Index
	suc.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/subnets", suc.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/subnets/{name}-{prefix}", suc.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/subnets/id/{id}", suc.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/subnets", suc.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/subnets/{name}-{prefix}", suc.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/subnets/{name}-{prefix}", suc.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/subnets/id/{id}", suc.RemoveByID).Methods("DELETE")

	// Image
	// Get Controller instance
	ic := controllers.NewImageController(getSession())

	// Set Database
	ic.SetDatabase(*database)

	// Set Schema URI
	ic.SetSchemaURI(*schemaURI)

	// Create Index
	ic.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/images", ic.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/images/{name}", ic.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/images/id/{id}", ic.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/images", ic.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/images/{name}", ic.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/images/{name}", ic.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/images/id/{id}", ic.RemoveByID).Methods("DELETE")

	// Image Version
	// Get Controller instance
	vc := controllers.NewImageVersionController(getSession())

	// Set Database
	vc.SetDatabase(*database)

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/image-versions", vc.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/image-versions/{name}", vc.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/image-versions/id/{id}", vc.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/image-versions", vc.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/image-versions/{name}", vc.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/image-versions/{name}", vc.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/image-versions/id/{id}", vc.RemoveByID).Methods("DELETE")

	// Boot Image
	// Get Controller instance
	bc := controllers.NewBootImageController(getSession())

	// Set Database
	bc.SetDatabase(*database)

	// Set Schema URI
	bc.SetSchemaURI(*schemaURI)

	// Create Index
	bc.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/boot-images", bc.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/boot-images/{name}", bc.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/boot-images/id/{id}", bc.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/boot-images", bc.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/boot-images/{name}", bc.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/boot-images/{name}", bc.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/boot-images/id/{id}", bc.RemoveByID).Methods("DELETE")

	// Boot Image Versions
	// Get Controller instance
	bcv := controllers.NewBootImageVersionController(getSession())

	// Set Database
	bcv.SetDatabase(*database)

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/boot-image-versions", bcv.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/boot-image-versions/{name}", bcv.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/boot-image-versions/id/{id}", bcv.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/boot-image-versions", bcv.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/boot-image-versions/{name}", bcv.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/boot-image-versions/{name}", bcv.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/boot-image-versions/id/{id}", bcv.RemoveByID).Methods("DELETE")

	// Tenant
	// Get Controller instance
	tc := controllers.NewTenantController(getSession())

	// Set Database
	tc.SetDatabase(*database)

	// Set Schema URI
	tc.SetSchemaURI(*schemaURI)

	// Create Index
	tc.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/tenants", tc.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/tenants/{name}", tc.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/tenants/id/{id}", tc.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/tenants", tc.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/tenants/{name}", tc.Update).Methods("PUT")
	r.HandleFunc("/"+version.APIVersion+"/tenants/{name}", tc.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/tenants/id/{id}", tc.RemoveByID).Methods("DELETE")

	// PXE Menu
	// Get Controller instance
	pc := controllers.NewPXEMenuController(getSession())

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/ipxe/{hwaddr}", pc.PXEMenu).Methods("GET")

	// Generic
	// Static files
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
