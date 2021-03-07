package post

import "fmt"

var (
	errIdLessThanZero = fmt.Errorf("id should be greated than zero")
	errNilPost        = fmt.Errorf("post can't be nil")
	errEmptyPostBody  = fmt.Errorf("post body can't be empty")
)
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
		return nil, errIdLessThanZero
	}

	return s.repo.PostsByUserId(id)
}

func (s *Service) Add(post *Post, userId int) error {
	if post == nil {
		return errNilPost
	}
	if userId <= 0 {
		return errIdLessThanZero
	}
	if post.Body == "" {
		return errEmptyPostBody
	}

	return s.repo.Add(post, userId)
}
