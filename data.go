package main

import "time"

type Card struct {
	ID           int           `json:"id"`
	Name         string        `json:"label"`
	StageID      int           `json:"stage"`
	Details      string        `json:"details"`
	StartDate    *time.Time    `json:"start_date"`
	OwnerID      int           `json:"owner"`
	Index        int           `json:"i"`
	AttachedData []*BinaryData `json:"attached,omitempty"`
}

type CardUpdate struct {
	Name         string        `json:"label"`
	StageID      FuzzyInt      `json:"stage"`
	Details      string        `json:"details"`
	StartDate    *time.Time    `json:"start_date"`
	OwnerID      FuzzyInt      `json:"owner"`
	AttachedData []*BinaryData `json:"attached,omitempty"`
}

type CardMove struct {
	Card   CardUpdate `json:"card"`
	Before FuzzyInt   `json:"before"`
}

type Stage struct {
	ID   int    `json:"id"`
	Name string `json:"label"`
}

type StageUpdate struct {
	Name string `json:"label"`
}

type BinaryData struct {
	ID     int    `json:"id"`
	Path   string `json:"-"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	CardID int    `json:"-"`
	Card   *Card  `json:"-"`
}
