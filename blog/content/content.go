package content

import (
	"time"
)

// Home represents an entry in the content index.
type Home struct {
	Slug    string    `json:"slug"`
	Date    time.Time `json:"date"`
	Summary string    `json:"summary"`
}

// contentRoot is the root directory for content files.
const contentRoot = "data/content"

// GetContentRoot returns the root directory for content files.
func GetContentRoot() string {
	return contentRoot
}
