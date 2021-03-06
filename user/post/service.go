package post

import "fmt"

type repository interface {
	PostsByUserId(id int) ([]Post, error)
	Add(post *Post, userId int) error
}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) PostsByUserId(id int) ([]Post, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id should be greated than zero")
	}

	return s.repo.PostsByUserId(id)
}

func (s *Service) Add(post *Post, userId int) error {
	if post == nil {
		return fmt.Errorf("post can't be nil")
	}
	if userId <= 0 {
		return fmt.Errorf("user id should be greater than zero")
	}
	if post.Body == "" {
		return fmt.Errorf("post body can't be empty")
	}

	return s.repo.Add(post, userId)
}
