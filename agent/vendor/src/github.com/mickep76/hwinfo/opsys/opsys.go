package opsys

// OpSys information.
type OpSys struct {
	Kernel         string `json:"kernel"`
	KernelVersion  string `json:"kernel_version"`
	Product        string `json:"product"`
	ProductVersion string `json:"product_version"`
}
