package models

import "gopkg.in/mgo.v2/bson"

type Image struct {
	ID               bson.ObjectId  `json:"id" bson:"_id"`
	Image            string         `json:"image" bson:"image"`
	Type             string         `json:"type" bson:"type"`
	BootImageID      bson.ObjectId  `json:"bootImageId" bson:"bootImageId"`
	BootImageRef     string         `json:"bootImageRef,omitempty"`
	BootImage        *BootImage     `json:"bootImage,omitempty"`
	BootImageVersion string         `json:"bootImageVersion" bson:"bootImageVersion"`
	Versions         []ImageVersion `json:"versions,omitempty" bson:"versions,omitempty"`
}

type ImageVersion struct {
	Version string `json:"version" bson:"version"`
	Created string `json:"created" bson:"created"`
}
