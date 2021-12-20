// Package note provides functions for interacting with notes
package note

// Note is a value object
type Note struct {
	ID   string `json:"id"`
	Note string `json:"note"`
}
