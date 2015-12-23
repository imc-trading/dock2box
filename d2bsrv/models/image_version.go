package models

import "gopkg.in/mgo.v2/bson"

type ImageVersion struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	ImageID  bson.ObjectId `json:"imageId" bson:"imageId"`
	ImageRef string        `json:"imageRef,omitempty"`
	Image    *Image        `json:"image,omitempty"`
	Version  string        `json:"version" bson:"version"`
	Created  string        `json:"created" bson:"created"`
}
