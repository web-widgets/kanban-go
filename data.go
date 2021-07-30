package main

type Card struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	StageID      int           `json:"stage"`
	AttachedData []*BinaryData `json:"attached,omitempty"`
}

type CardUpdate struct {
	Name         string        `json:"name"`
	StageID      int           `json:"stage"`
	AttachedData []*BinaryData `json:"attached,omitempty"`
}

type Stage struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type StageUpdate struct {
	Name string `json:"name"`
}

type BinaryData struct {
	ID     int    `json:"id"`
	Path   string `json:"-"`
	Name   string `json:"name"`
	CardID int    `json:"-"`
	Card   *Card  `json:"-"`
}
