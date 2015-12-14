package main

// TODO
// - Add filters for /subnets

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
	flag.Parse()

	// Print version.
	if *appVersion {
		fmt.Printf("d2bsrv %s\n", version.Version)
		os.Exit(0)
	}

	// Create new router
	r := mux.NewRouter()

	// Host
	// Get Controller instance
	hc := controllers.NewHostController(getSession())

	// Create Index
	hc.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/hosts", hc.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/hosts/{name}", hc.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/hosts/id/{id}", hc.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/hosts", hc.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/hosts/{name}", hc.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/hosts/id/{id}", hc.RemoveByID).Methods("DELETE")

	// Site
	// Get Controller instance
	sc := controllers.NewSiteController(getSession())

	// Create Index
	sc.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/sites", sc.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/sites/{name}", sc.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/sites/id/{id}", sc.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/sites", sc.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/sites/{name}", sc.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/sites/id/{id}", sc.RemoveByID).Methods("DELETE")

	// Subnet
	// Get Controller instance
	suc := controllers.NewSubnetController(getSession())

	// Create Index
	suc.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/subnets", suc.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/subnets/{name}-{prefix}", suc.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/subnets/id/{id}", suc.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/subnets", suc.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/subnets/{name}/{prefix}", suc.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/subnets/id/{id}", suc.RemoveByID).Methods("DELETE")

	// Image
	// Get Controller instance
	ic := controllers.NewImageController(getSession())

	// Create Index
	ic.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/images", ic.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/images/{name}", ic.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/images/id/{id}", ic.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/images", ic.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/images/{name}", ic.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/images/id/{id}", ic.RemoveByID).Methods("DELETE")

	// Image Versions
	// Get Controller instance
	vc := controllers.NewImageVersionController(getSession())

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/images/{name}/versions", vc.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/images/id/{id}/versions", vc.AllByID).Methods("GET")

	// Boot Image
	// Get Controller instance
	bc := controllers.NewBootImageController(getSession())

	// Create Index
	bc.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/boot-images", bc.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/boot-images/{name}", bc.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/boot-images/id/{id}", bc.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/boot-images", bc.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/boot-images/{name}", bc.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/boot-images/id/{id}", bc.RemoveByID).Methods("DELETE")

	// Tenant
	// Get Controller instance
	tc := controllers.NewTenantController(getSession())

	// Create Index
	tc.CreateIndex()

	// Add handlers for endpoints
	r.HandleFunc("/"+version.APIVersion+"/tenants", tc.All).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/tenants/{name}", tc.Get).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/tenants/id/{id}", tc.GetByID).Methods("GET")
	r.HandleFunc("/"+version.APIVersion+"/tenants", tc.Create).Methods("POST")
	r.HandleFunc("/"+version.APIVersion+"/tenants/{name}", tc.Remove).Methods("DELETE")
	r.HandleFunc("/"+version.APIVersion+"/tenants/id/{id}", tc.RemoveByID).Methods("DELETE")

	// Generic
	// Static files
	schemas := http.StripPrefix("/"+version.APIVersion+"/schemas/", http.FileServer(http.Dir("schemas")))
	r.PathPrefix("/" + version.APIVersion + "/schemas/").Handler(schemas)
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
