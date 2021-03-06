package post

import (
	"time"
)

type Post struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Body      string
}
