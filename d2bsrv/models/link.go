package models

type Link struct {
	HRef   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}
