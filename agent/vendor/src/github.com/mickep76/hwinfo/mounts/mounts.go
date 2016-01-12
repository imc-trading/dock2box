package mounts

type Mount struct {
	Source  string `json:"source"`
	Target  string `json:"target"`
	FSType  string `json:"fs_type"`
	Options string `json:"options"`
}
