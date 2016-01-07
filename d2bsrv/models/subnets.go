package models

import "gopkg.in/mgo.v2/bson"

type Subnet struct {
	ID     bson.ObjectId `field:"id" json:"id,omitempty" bson:"_id"`
	Subnet string        `field:"subnet" json:"subnet,omitempty" bson:"subnet"`
	Mask   string        `field:"mask" json:"mask,omitempty" bson:"mask"`
	Gw     string        `field:"gw" json:"gw,omitempty" bson:"gw"`
	SiteID bson.ObjectId `field:"siteId" json:"siteId,omitempty" bson:"siteId"`
	Site   *Site         `json:"site,omitempty"`
}
