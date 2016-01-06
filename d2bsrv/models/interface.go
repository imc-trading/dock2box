package models

import "gopkg.in/mgo.v2/bson"

type Interface struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Interface string        `json:"interface" bson:"interface"`
	DHCP      bool          `json:"dhcp" bson:"dhcp"`
	IPv4      string        `json:"ipv4,omitempty" bson:"ipv4,omitempty"`
	HwAddr    string        `json:"hwAddr" bson:"hwAddr"`
	SubnetID  bson.ObjectId `json:"subnetId,omitempty" bson:"subnetId,omitempty"`
	Subnet    *Subnet       `json:"subnet,omitempty"`
	HostID    bson.ObjectId `json:"hostId" bson:"hostId"`
}
