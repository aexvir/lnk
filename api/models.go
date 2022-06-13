package api

import _ "embed"

type Link struct {
	Slug      string            `json:"slug"`
	Target    string            `json:"target"`
	Hits      uint64            `json:"hits"`
	Histogram map[string]uint64 `json:"histogram"`
}

type LinkResp struct {
	Link string `json:"link"`
}

type CreateLinkReq struct {
	Target string  `json:"target"`
	Slug   *string `json:"slug"`
}

//go:embed docs/openapi.json
var Schema []byte
