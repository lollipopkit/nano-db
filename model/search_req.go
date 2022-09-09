package model

type SearchReq struct {
	Path  string `json:"path"`
	Regex string `json:"regex"`
}
