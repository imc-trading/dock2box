package models

import "gopkg.in/mgo.v2/bson"

type ImageTag struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Tag      string        `json:"version" bson:"version"`
	Created  string        `json:"created" bson:"created"`
	ImageID  bson.ObjectId `json:"imageId" bson:"imageId"`
	ImageRef string        `json:"imageRef,omitempty"`
	Image    *Image        `json:"image,omitempty"`
}
