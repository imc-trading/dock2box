package controllers

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/imc-trading/dock2box/d2bsrv/models"
)

func Center(size int, deco string, str string) string {
	if size < len(str) {
		return str
	}

	pad := (size - len(str)) / 2
	lpad := pad
	rpad := size - len(str) - lpad

	return fmt.Sprintf("%s%s%s", strings.Repeat(deco, lpad), str, strings.Repeat(deco, rpad))
}

var funcs = template.FuncMap{
	"center": Center,
}

var templates = template.Must(template.New("main").Funcs(funcs).ParseGlob("templates/*.tmpl"))

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
	ImageTags   []models.ImageTag
	Host        models.Host
	Subnet      models.Subnet
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
		http.Error(w, "Can't get images", http.StatusInternalServerError)
		return
	}

	// Get images tags and embed them inside images.
	for i, e := range input.Images {
		tags := make([]models.ImageTag, 0)
		input.Images[i].Tags = &tags
		if err := c.session.DB(c.database).C("image_tags").Find(bson.M{"imageId": e.ID}).All(input.Images[i].Tags); err != nil {
			http.Error(w, fmt.Sprintf("Can't get image tags for image id: %s", e.ID), http.StatusInternalServerError)
			return
		}
	}

	// Get host.
	if err := c.session.DB(c.database).C("hosts").Find(bson.M{"interfaces": bson.M{"$elemMatch": bson.M{"hwAddr": input.HWAddr}}}).One(&input.Host); err != nil {

		// Unregistered host, get subnet.
		if err := c.session.DB(c.database).C("subnets").Find(bson.M{"subnet": fmt.Sprintf("%s/%d", input.Network, input.Prefix)}).One(&input.Subnet); err != nil {
			// Print debug.
			if input.Debug == "true" {
				jsonWriter(w, r, input, http.StatusOK)
				return
			}

			// Template menu.
			if err := templates.ExecuteTemplate(w, "night", input); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			if err := templates.ExecuteTemplate(w, "no_subnet", input); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Unregistered host, get site.
		if err := c.session.DB(c.database).C("sites").Find(bson.M{"_id": input.Subnet.SiteID}).One(&input.Subnet.Site); err != nil {
			http.Error(w, fmt.Sprintf("Can't get site specified for subnet: %s/%d", input.Network, input.Prefix), http.StatusInternalServerError)
			return
		}

		// Print debug.
		if input.Debug == "true" {
			jsonWriter(w, r, input, http.StatusOK)
			return
		}

		// Template menu.
		if err := templates.ExecuteTemplate(w, input.Subnet.Site.PXETheme, input); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := templates.ExecuteTemplate(w, "unregistered", input); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Get image tag.
	if err := c.session.DB(c.database).C("image_tags").FindId(input.Host.ImageTagID).One(&input.Host.ImageTag); err != nil {
		http.Error(w, fmt.Sprintf("Can't get image tag id: %s", input.Host.ImageTagID.Hex()), http.StatusInternalServerError)
		return
	}

	// Get image.
	if err := c.session.DB(c.database).C("images").FindId(input.Host.ImageTag.ImageID).One(&input.Host.Image); err != nil {
		http.Error(w, fmt.Sprintf("Can't get image id: %s", input.Host.ImageTag.ImageID.Hex()), http.StatusInternalServerError)
		return
	}

	// Get boot image tag.
	if err := c.session.DB(c.database).C("image_tags").FindId(input.Host.Image.BootImageTagID).One(&input.Host.BootImageTag); err != nil {
		http.Error(w, fmt.Sprintf("Can't get boot image tag id: %s", input.Host.Image.BootImageTagID.Hex()), http.StatusInternalServerError)
		return
	}

	// Get boot image.
	if err := c.session.DB(c.database).C("images").FindId(input.Host.BootImageTag.ImageID).One(&input.Host.BootImage); err != nil {
		http.Error(w, fmt.Sprintf("Can't get boot image id: %s", input.Host.BootImageTag.ImageID.Hex()), http.StatusInternalServerError)
		return
	}

	// Get tenant.
	if err := c.session.DB(c.database).C("tenants").FindId(input.Host.TenantID).One(&input.Host.Tenant); err != nil {
		http.Error(w, fmt.Sprintf("Can't get tenant id: %s", input.Host.TenantID.Hex()), http.StatusInternalServerError)
		return
	}

	// Get site.
	if err := c.session.DB(c.database).C("sites").FindId(input.Host.SiteID).One(&input.Host.Site); err != nil {
		http.Error(w, fmt.Sprintf("Can't get site id: %s", input.Host.SiteID.Hex()), http.StatusInternalServerError)
		return
	}

	// Print debug.
	if input.Debug == "true" {
		jsonWriter(w, r, input, http.StatusOK)
		return
	}

	// Template menu.
	if err := templates.ExecuteTemplate(w, input.Host.Site.PXETheme, input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := templates.ExecuteTemplate(w, "registered", input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
