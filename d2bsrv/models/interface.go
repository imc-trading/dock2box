package models

import "gopkg.in/mgo.v2/bson"

type Interface struct {
	ID        bson.ObjectId `field:"id" json:"id" bson:"_id"`
	Interface string        `field:"interface" json:"interface" bson:"interface"`
	DHCP      bool          `field:"dhcp" json:"dhcp" bson:"dhcp"`
	IPv4      string        `field:"ipv4" json:"ipv4,omitempty" bson:"ipv4,omitempty"`
	HwAddr    string        `field:"hwAddr" json:"hwAddr" bson:"hwAddr"`
	SubnetID  bson.ObjectId `field:"subnetId" json:"subnetId,omitempty" bson:"subnetId,omitempty"`
	Subnet    *Subnet       `json:"subnet,omitempty"`
	HostID    bson.ObjectId `field:"hostId" json:"hostId" bson:"hostId"`
	Host      *Host         `json:"host,omitempty"`
	Links     *[]Link       `json:"links,omitempty"`
}
