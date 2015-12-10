package models

import "gopkg.in/mgo.v2/bson"

type Tenant struct {
	ID     bson.ObjectId `json:"id" bson:"_id"`
	Tenant string        `json:"tenant" bson:"tenant"`
}
