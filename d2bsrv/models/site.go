package models

import "gopkg.in/mgo.v2/bson"

type Site struct {
	ID                 bson.ObjectId `json:"id" bson:"_id"`
	Site               string        `json:"site" bson:"site"`
	Domain             string        `json:"domain" bson:"domain"`
	DNS                []string      `json:"dns" bson:"dns"`
	DockerRegistry     string        `json:"dockerRegistry" bson:"dockerRegistry"`
	ArtifactRepository string        `json:"artifactRepository" bson:"artifactRepository"`
	PXETheme           string        `json:"pxeTheme" bson:"pxeTheme"`
	NamingScheme       string        `json:"namingScheme" bson:"namingScheme"`
}
