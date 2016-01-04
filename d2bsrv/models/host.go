package models

import "gopkg.in/mgo.v2/bson"

type Host struct {
	ID              bson.ObjectId   `json:"id" bson:"_id"`
	Host            string          `json:"host" bson:"host"`
	Build           bool            `json:"build" bson:"build"`
	Debug           bool            `json:"debug" bson:"debug"`
	GPT             bool            `json:"gpt" bson:"gpt"`
	ImageVersionID  bson.ObjectId   `json:"imageVersionId" bson:"imageVersionId"`
	ImageVersionRef string          `json:"imageVersionRef" bson:"imageVersionRef"`
	ImageVersion    *ImageVersion   `json:"imageVersion,omitempty"`
	KOpts           string          `json:"kOpts" bson:"kOpts"`
	TenantID        bson.ObjectId   `json:"tenantId" bson:"tenantId"`
	TenantRef       string          `json:"tenantRef" bson:"tenantRef"`
	Tenant          *Tenant         `json:"tenant,omitempty"`
	Labels          []string        `json:"labels" bson:"labels"`
	SiteID          bson.ObjectId   `json:"siteId" bson:"siteId"`
	SiteRef         string          `json:"siteRef" bson:"siteRef"`
	Site            *Site           `json:"site,omitempty"`
	Interfaces      []HostInterface `json:"interfaces" bson:"interfaces"`
}

type HostInterface struct {
	Interface string        `json:"interface" bson:"interface"`
	DHCP      bool          `json:"dhcp" bson:"dhcp"`
	IPv4      string        `json:"ipv4,omitempty" bson:"ipv4,omitempty"`
	HwAddr    string        `json:"hwAddr" bson:"hwAddr"`
	SubnetID  bson.ObjectId `json:"subnetId,omitempty" bson:"subnetId,omitempty"`
	SubnetRef string        `json:"subnetRef,omitempty" bson:"subnetRef,omitempty"`
	Subnet    *Subnet       `json:"subnet,omitempty"`
}
