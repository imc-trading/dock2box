package models

import "gopkg.in/mgo.v2/bson"

type Tenant struct {
	ID     bson.ObjectId `field:"id" json:"id,omitempty" bson:"_id"`
	Tenant string        `field:"tenant" json:"tenant,omitempty" bson:"tenant"`
	Links  *[]Link       `json:"links,omitempty"`
}
