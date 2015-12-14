package models

import "gopkg.in/mgo.v2/bson"

type BootImage struct {
	ID       bson.ObjectId      `json:"id" bson:"_id"`
	Image    string             `json:"image" bson:"image"`
	KOpts    string             `json:"kOpts" bson:"kOpts"`
	Versions []BootImageVersion `json:"versions,omitempty" bson:"versions,omitempty"`
}

type BootImageVersion struct {
	Version string `json:"version" bson:"version"`
	Created string `json:"created" bson:"created"`
}
