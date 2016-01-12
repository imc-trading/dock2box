package command

import (
	"log"

	"github.com/imc-trading/dock2box/cli/prompt"
	"github.com/imc-trading/dock2box/client"
)

func chooseImage(clnt *client.Client, imageID string) *string {
	r, err := clnt.Image.All()
	if err != nil {
		log.Fatalf(err.Error())
	}

	images := *r
	var list []string
	def := -1
	for i, v := range images {
		if v.ID == imageID {
			def = i
		}
		list = append(list, v.Image+" type: "+v.Type)
	}
	return &images[prompt.Choice("Choose image", def, list)].ID
}

func chooseTag(clnt *client.Client, tagID string) *string {
	chooseImage(clnt, "")
	// imageID := chooseImage(clnt, "")
	// r, err := clnt.Tag.Match("imageId", imageID)
	r, err := clnt.Tag.All()
	if err != nil {
		log.Fatalf(err.Error())
	}

	tags := *r
	var list []string
	def := -1
	for i, v := range tags {
		if v.ID == tagID {
			def = i
		}
		list = append(list, v.Tag) // Same as in template
	}
	return &tags[prompt.Choice("Choose tag", def, list)].ID
}

func chooseTenants(clnt *client.Client, tenantID string) *string {
	r, err := clnt.Tenant.All()
	if err != nil {
		log.Fatalf(err.Error())
	}

	tenants := *r
	var list []string
	def := -1
	for i, v := range tenants {
		if v.ID == tenantID {
			def = i
		}
		list = append(list, v.Tenant)
	}
	return &tenants[prompt.Choice("Choose tenant", def, list)].ID
}

func chooseSite(clnt *client.Client, siteID string) *string {
	r, err := clnt.Site.All()
	if err != nil {
		log.Fatalf(err.Error())
	}

	sites := *r
	var list []string
	def := -1
	for i, v := range sites {
		if v.ID == siteID {
			def = i
		}
		list = append(list, v.Site+", domain: "+v.Domain)
	}
	return &sites[prompt.Choice("Choose site", def, list)].ID
}

func chooseSubnet(clnt *client.Client, siteID string, subnetID string) *string {
	r, err := clnt.Subnet.All()
	if err != nil {
		log.Fatalf(err.Error())
	}

	subnets := *r
	var list []string
	def := -1
	for i, v := range subnets {
		// UGLY: keep until backend supports filters
		if v.SiteID == siteID {
			if v.ID == subnetID {
				def = i
			}
			list = append(list, v.Subnet)
		}
	}
	return &subnets[prompt.Choice("Choose subnet", def, list)].ID
}
