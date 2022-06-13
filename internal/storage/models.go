package storage

type Link struct {
	Slug      string            `json:"slug"`
	Target    string            `json:"target"`
	Hits      uint64            `json:"hits"`
	Histogram map[string]uint64 `json:"histogram"`
}
