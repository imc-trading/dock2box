package routes

// Info structure for system routes.
type Route struct {
	Destination string `json:"destination"`
	Gateway     string `json:"gateway"`
	Genmask     string `json:"genmask"`
	Flags       string `json:"flags"`
	MSS         int    `json:"mss"` // Maximum segment size
	Window      int    `json:"window"`
	IRTT        int    `json:"irtt"` // Initial round trip time
	Interface   string `json:"interface"`
}
