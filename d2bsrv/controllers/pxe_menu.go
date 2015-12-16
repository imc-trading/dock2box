package controllers

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/imc-trading/dock2box/d2bsrv/models"
)

var templates = template.Must(template.ParseFiles("templates/menu.html"))

type Input struct {
	HWAddr      string
	IPv4        string
	Netmask     string
	Gateway     string
	Prefix      int
	Network     string
	Serial      string
	BoardSerial string
	Debug       string
	Images      []models.Image
	Host        models.Host
}

type PXEMenuController struct {
	database string
	session  *mgo.Session
}

func NewPXEMenuController(s *mgo.Session) *PXEMenuController {
	return &PXEMenuController{
		database: "d2b",
		session:  s,
	}
}

func (c PXEMenuController) SetDatabase(database string) {
	c.database = database
}

func (c PXEMenuController) PXEMenu(w http.ResponseWriter, r *http.Request) {
	input := Input{}

	// Get and parse hwaddr.
	input.HWAddr = mux.Vars(r)["hwaddr"]
	if _, err := net.ParseMAC(input.HWAddr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get and parse IP.
	input.IPv4 = r.URL.Query().Get("ipv4")
	if net.ParseIP(input.IPv4) == nil {
		http.Error(w, fmt.Sprintf("invalid IP %s", input.IPv4), http.StatusBadRequest)
		return
	}

	// Get and parse netmask.
	input.Netmask = r.URL.Query().Get("netmask")
	maskIP := net.ParseIP(input.Netmask).To4()
	if maskIP == nil {
		http.Error(w, fmt.Sprintf("invalid Netmask %s", input.Netmask), http.StatusBadRequest)
		return
	}

	// Get netmask prefix.
	mask := net.IPMask(maskIP)
	input.Prefix, _ = mask.Size()

	// Get and parse network.
	_, network, err := net.ParseCIDR(input.IPv4 + "/" + strconv.Itoa(input.Prefix))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	input.Network = network.IP.String()

	// Get gateway.
	input.Gateway = r.URL.Query().Get("gateway")
	if net.ParseIP(input.Gateway) == nil {
		http.Error(w, fmt.Sprintf("invalid Gateway %s", input.Gateway), http.StatusBadRequest)
		return
	}

	// Optional.
	input.Serial = r.URL.Query().Get("serial")
	input.BoardSerial = r.URL.Query().Get("boardserial")
	input.Debug = r.URL.Query().Get("debug")

	// Get all images.
	if err := c.session.DB(c.database).C("images").Find(nil).All(&input.Images); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get host.
	if err := c.session.DB(c.database).C("hosts").Find(bson.M{"interfaces": bson.M{"$elemMatch": bson.M{"hwAddr": input.HWAddr}}}).One(&input.Host); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get image.
	if err := c.session.DB(c.database).C("images").FindId(input.Host.ImageID).One(&input.Host.Image); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get boot image.
	if err := c.session.DB(c.database).C("boot_images").FindId(input.Host.Image.BootImageID).One(&input.Host.Image.BootImage); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get tenant.
	if err := c.session.DB(c.database).C("tenants").FindId(input.Host.TenantID).One(&input.Host.Tenant); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get site.
	if err := c.session.DB(c.database).C("sites").FindId(input.Host.SiteID).One(&input.Host.Site); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Print debug.
	if input.Debug == "true" {
		jsonWriter(w, r, input, http.StatusOK)
		return
	}

	// Template menu.
	templates.ExecuteTemplate(w, "menu", input)
}
