// Package rounds contains data and calculations independent of HTTP and SQL.
package rounds

// HoleScore describes one hole. JSON tags map Go fields to request field names.
type HoleScore struct {
	HoleNumber int `json:"hole_number"`
	Par        int `json:"par"`
	Strokes    int `json:"strokes"`
}

// Round describes a completed round. Date is an RFC 3339 timestamp with a timezone.
type Round struct {
	Username   string      `json:"username"`
	Id         int64       `json:"id"`
	CourseName string      `json:"course_name"`
	Date       string      `json:"date"`
	Holes      []HoleScore `json:"holes"`
}

// ValidateRound checks these fields before HTTP calculations.

// User identifies a player stored in the users table.
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}
