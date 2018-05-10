package model

import "time"

type Site struct {
	UUID    string     `json:"uuid"`
	Created time.Time  `json:"created"`
	Updated *time.Time `json:"updated,omitempty"`
	Name    string     `json:"name"`
	Descr   string     `json:"descr,omitempty"`
	Region  string     `json:"region,omitempty"`
	Country string     `json:"country,omitempty"`
	Lat     float32    `json:"lat,omitempty"`
	Lng     float32    `json:"lng,omitempty"`
}

type Sites []*Site
