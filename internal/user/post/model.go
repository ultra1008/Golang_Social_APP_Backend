package post

import (
	"encoding/json"
	"sort"
	"time"
)

type Author struct {
	ID        int
	FirstName string
	LastName  string
	Login     string
}
type Post struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Body      string
	Author    Author
}

func (p Post) AsByteJSON() ([]byte, error) {
	return json.Marshal(p)
}

type Feed []Post

func (f Feed) Sort() {
	sort.Slice(f, func(i, j int) bool {
		return f[i].CreatedAt.After(f[j].CreatedAt)
	})
}
