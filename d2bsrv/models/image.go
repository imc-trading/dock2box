package models

import "gopkg.in/mgo.v2/bson"

type Image struct {
	ID              bson.ObjectId `json:"id" bson:"_id"`
	Image           string        `json:"image" bson:"image"`
	SHA256          string        `json:"sha256" bson:"sha256"`
	Type            string        `json:"type" bson:"type"`
	KOpts           string        `json:"kOpts,omitempty" bson:"kOpts,omitempty"`
	BootImageTagID  bson.ObjectId `json:"bootImageTagId,omitempty" bson:"bootImageTagId,omitempty"`
	BootImageTagRef string        `json:"bootImageTagRef,omitempty" bson:"bootImageTagRef,omitempty"`
	BootImageTag    *ImageTag     `json:"bootImageTag,omitempty"`
	BootImage       *Image        `json:"bootImage,omitempty"`
}
