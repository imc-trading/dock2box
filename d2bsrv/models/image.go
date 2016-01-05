package models

import "gopkg.in/mgo.v2/bson"

type Image struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	Image      string        `json:"image" bson:"image"`
	Type       string        `json:"type" bson:"type"`
	KOpts      string        `json:"kOpts,omitempty" bson:"kOpts,omitempty"`
	BootTagID  bson.ObjectId `json:"bootTagId,omitempty" bson:"bootTagId,omitempty"`
	BootTagRef string        `json:"bootTagRef,omitempty" bson:"bootTagRef,omitempty"`
	BootTag    *Tag          `json:"bootTag,omitempty"`
	BootImage  *Image        `json:"bootImage,omitempty"`
	Tags       *[]Tag        `json:"tags,omitempty"`
}
