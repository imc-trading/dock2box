package lvm

type LVM struct {
	PhysVols *[]PhysVol `json:"phys_vols"`
	LogVols  *[]LogVol  `json:"log_vols"`
	VolGrps  *[]VolGrp  `json:"vol_grps"`
}

type PhysVol struct {
	Name   string `json:"name"`
	VolGrp string `json:"vol_group"`
	Format string `json:"format"`
	Attr   string `json:"attr"`
	SizeGB int    `json:"size_gb"`
	FreeGB int    `json:"free_gb"`
}

type LogVol struct {
	Name   string `json:"name"`
	VolGrp string `json:"vol_grp"`
	Attr   string `json:"attr"`
	SizeGB int    `json:"size_gb"`
}

type VolGrp struct {
	Name   string `json:"name"`
	Attr   string `json:"attr"`
	SizeGB int    `json:"size_gb"`
	FreeGB int    `json:"free_gb"`
}
