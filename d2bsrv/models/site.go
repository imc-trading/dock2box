package models

import "gopkg.in/mgo.v2/bson"

type Site struct {
	ID                 bson.ObjectId `field:"id" json:"id,omitempty" bson:"_id"`
	Site               string        `field:"site" json:"site,omitempty" bson:"site"`
	Domain             string        `field:"domain" json:"domain,omitempty" bson:"domain"`
	DNS                []string      `field:"dns" json:"dns,omitempty" bson:"dns"`
	DockerRegistry     string        `field:"dockerRegistry" json:"dockerRegistry,omitempty" bson:"dockerRegistry"`
	ArtifactRepository string        `field:"artifactRepository" json:"artifactRepository,omitempty" bson:"artifactRepository"`
	PXETheme           string        `field:"pxeTheme" json:"pxeTheme,omitempty" bson:"pxeTheme"`
	NamingScheme       string        `field:"namingScheme" json:"namingScheme,omitempty" bson:"namingScheme"`
}
