package post

import (
	"fmt"

	"github.com/niklod/highload-social-network/internal/cache"
)

var (
	errIdLessThanZero = fmt.Errorf("id should be greated than zero")
	errNilPost        = fmt.Errorf("post can't be nil")
	errEmptyPostBody  = fmt.Errorf("post body can't be empty")
)

type repository interface {
	PostsByUserId(id int) ([]Post, error)
	UserFeed(id int) (Feed, error)
	Add(post *Post, userId int) error
}

type Service struct {
	repo  repository
	cache cache.Cache
}

func NewService(repo repository, cache cache.Cache) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
	}
}

func (s *Service) UserFeed(userId int) (Feed, error) {
	if userId <= 0 {
		return nil, errIdLessThanZero
	}

	v, ok := s.cache.Read(userId)
	if ok {
		userFeed, ok := v.(Feed)
		if !ok {
			return nil, fmt.Errorf("post.Service: %v", cache.ErrInvalidCacheItem)
		}

		return userFeed, nil
	}

	userFeed, err := s.repo.UserFeed(userId)
	if err != nil {
		return nil, fmt.Errorf("post.Service: %v", err)
	}

	go s.cache.Write(userId, userFeed)

	return userFeed, nil
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

	err := s.repo.Add(post, userId)
	if err != nil {
		return fmt.Errorf("post.Service: %v", err)
	}

	// add to queue

	return nil
}
