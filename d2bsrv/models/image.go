package models

import "gopkg.in/mgo.v2/bson"

type Image struct {
	ID                  bson.ObjectId `json:"id" bson:"_id"`
	Image               string        `json:"image" bson:"image"`
	Type                string        `json:"type" bson:"type"`
	KOpts               string        `json:"kOpts,omitempty" bson:"kOpts,omitempty"`
	BootImageVersionID  bson.ObjectId `json:"bootImageVersionId,omitempty" bson:"bootImageVersionId,omitempty"`
	BootImageVersionRef string        `json:"bootImageVersionRef,omitempty" bson:"bootImageVersionRef,omitempty"`
	BootImageVersion    *ImageVersion `json:"bootImageVersion,omitempty"`
}
