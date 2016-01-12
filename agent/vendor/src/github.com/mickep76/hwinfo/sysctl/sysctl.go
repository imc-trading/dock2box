package sysctl

// Sysctl structure for sysctl key/values.
type Sysctl struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
