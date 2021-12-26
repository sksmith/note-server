// Package note provides functions for interacting with notes
package note

import (
	"time"
)

// A note as created by a user
type Note struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Data    string    `json:"data"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

// A note as represented in the index
type ListNote struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}
