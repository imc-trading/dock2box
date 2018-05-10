package model

import "time"

type Host struct {
	UUID           string         `json:"uuid"`
	Created        time.Time      `json:"created"`
	Updated        *time.Time     `json:"updated,omitempty"`
	ClientUUID     string         `json:"clientUUID"`
	Name           string         `json:"name"`
	IPAddr         string         `json:"ipAddr"`
	Hardware       *Hardware      `json:"hardware,omitempty"`
	HwAddr         string         `json:"hwAddr,omitempty"`
	SerialNumber   string         `json:"serialNumber,omitempty"`
	Role           *Role          `json:"role,omitempty"`
	Pool           *Pool          `json:"pool,omitempty"`
	Site           *Site          `json:"site,omitempty"`
	Rack           *Rack          `json:"rack,omitempty"`
	RackUnitPos    int            `json:"rackUnitPos,omitempty"`
	RackUnitHeight int            `json:"rackUnitHeight,omitempty"`
	Tenant         *Tenant        `json:"tenantUUID,omitempty"`
	DockerImage    *DockerImage   `json:"dockerImage,omitempty"`
	AllowBuild     bool           `json:"allowBuild"`
	SecurityGroup  *SecurityGroup `json:"securityGroup"`
	HostGroups     HostGroups     `json:"hostGroups,omitempty"`
}

type Hosts []*Host
