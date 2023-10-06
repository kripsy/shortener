package models

// Stats represent an information about data in service
// count urls and users.
type Stats struct {
	URLs  int `json:"urls,omitempty"`
	Users int `json:"users,omitempty"`
}
