package models

import "gopkg.in/mgo.v2/bson"

type Tag struct {
	ID      bson.ObjectId `field:"id" json:"id,omitempty" bson:"_id" `
	Tag     string        `field:"tag" json:"tag,omitempty" bson:"tag"`
	Created string        `field:"created" json:"created,omitempty" bson:"created"`
	SHA256  string        `field:"sha256" json:"sha256,omitempty" bson:"sha256"`
	ImageID bson.ObjectId `field:"imageId" json:"imageId,omitempty" bson:"imageId"`
	Image   *Image        `json:"image,omitempty"`
	Links   *[]Link       `json:"links,omitempty"`
}
