package post

import (
	"time"
)

type Feed []Post

type Post struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Body      string
	Author    Author
}

type Author struct {
	FirstName string
	LastName  string
	Login     string
}
