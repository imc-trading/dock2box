package models

import "gopkg.in/mgo.v2/bson"

type Subnet struct {
	ID     bson.ObjectId `json:"id" bson:"_id"`
	Subnet string        `json:"subnet" bson:"subnet"`
	Mask   string        `json:"mask" bson:"mask"`
	Gw     string        `json:"gw" bson:"gw"`
	SiteID bson.ObjectId `json:"siteId" bson:"siteId"`
	Site   *Site         `json:"site,omitempty"`
}
