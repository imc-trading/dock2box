package models

import "gopkg.in/mgo.v2/bson"

type Tag struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	Tag     string        `json:"tag" bson:"tag"`
	Created string        `json:"created" bson:"created"`
	SHA256  string        `json:"sha256" bson:"sha256"`
	ImageID bson.ObjectId `json:"imageId" bson:"imageId"`
	Image   *Image        `json:"image,omitempty"`
	Links   *[]Link       `json:"links,omitempty"`
}
