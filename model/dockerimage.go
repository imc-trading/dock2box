package model

import "time"

type DockerImage struct {
	UUID     string     `json:"uuid"`
	Created  time.Time  `json:"created"`
	Updated  *time.Time `json:"updated,omitempty"`
	Registry string     `json:"registry,omitempty"`
	Repo     string     `json:"repo,omitempty"`
	Name     string     `json:"name"`
	Tag      string     `json:"tag,omitempty"`
}

type DockerImages []*DockerImage
