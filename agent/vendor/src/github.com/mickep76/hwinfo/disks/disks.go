package disks

// Disk information.
type Disk struct {
	Device string `json:"device"`
	Name   string `json:"name"`
	SizeGB int    `json:"size_gb"`
}
