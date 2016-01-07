package controllers

import (
	"net/http"

	"github.com/imc-trading/dock2box/d2bsrv/models"
)

type RootController struct {
	baseURI  string
	envelope bool
	hateoas  bool
}

func NewRootController(b string, e bool, h bool) *RootController {
	return &RootController{
		baseURI:  b,
		envelope: e,
		hateoas:  h,
	}
}

func (c *RootController) All(w http.ResponseWriter, r *http.Request) {
	// Initialize empty struct list
	s := models.Root{}

	// HATEOAS Links
	if c.hateoas == true || r.URL.Query().Get("hateoas") == "true" {
		s.Links = &[]models.Link{
			models.Link{
				HRef:   c.baseURI + "/hosts",
				Rel:    "self",
				Method: "GET",
			},
			models.Link{
				HRef:   c.baseURI + "/interfaces",
				Rel:    "self",
				Method: "GET",
			},
			models.Link{
				HRef:   c.baseURI + "/images",
				Rel:    "self",
				Method: "GET",
			},
			models.Link{
				HRef:   c.baseURI + "/tags",
				Rel:    "self",
				Method: "GET",
			},
			models.Link{
				HRef:   c.baseURI + "/sites",
				Rel:    "self",
				Method: "GET",
			},
			models.Link{
				HRef:   c.baseURI + "/subnets",
				Rel:    "self",
				Method: "GET",
			},
			models.Link{
				HRef:   c.baseURI + "/tenants",
				Rel:    "self",
				Method: "GET",
			},
		}
	}

	// Write content-type, header and payload
	jsonWriter(w, r, s, http.StatusOK, c.envelope)
}
