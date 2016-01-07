package models

import "gopkg.in/mgo.v2/bson"

type Image struct {
	ID        bson.ObjectId `field:"id" json:"id,omitempty" bson:"_id"`
	Image     string        `field:"image" json:"image,omitempty" bson:"image"`
	Type      string        `field:"type" json:"type,omitempty" bson:"type"`
	KOpts     string        `field:"kOpts" json:"kOpts,omitempty" bson:"kOpts,omitempty"`
	BootTagID bson.ObjectId `field:"bootTagId" json:"bootTagId,omitempty" bson:"bootTagId,omitempty"`
	BootTag   *Tag          `json:"bootTag,omitempty"`
	BootImage *Image        `json:"bootImage,omitempty"`
	Tags      *[]Tag        `json:"tags,omitempty"`
}
