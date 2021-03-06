package post

import (
	"database/sql"
	"fmt"
	"log"
)

type mysql struct {
	db *sql.DB
}

func NewRepository(client *sql.DB) repository {
	return &mysql{
		db: client,
	}
}

func (m *mysql) PostsByUserId(id int) ([]Post, error) {
	var posts []Post

	query, ctx, cancel := GetQuery(PostsByUserId)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("posts.PostByUserId - sending query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		post := Post{}

		err := rows.Scan(
			&post.ID,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Body,
		)
		if err != nil {
			log.Printf("posts.postbyuserid - scanning user: %v", err)
			continue
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("posts.postbyuserid - iterating through rows: %v", err)
	}

	return posts, nil
}

func (m *mysql) Add(post *Post, userId int) error {
	query, ctx, cancel := GetQuery(InsertPost)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, userId, post.Body)
	if err != nil {
		return fmt.Errorf("posts.Add - sending query: %v", err)
	}

	return nil
}
