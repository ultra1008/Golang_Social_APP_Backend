package post

import (
	"fmt"
	"time"

	"github.com/niklod/highload-social-network/internal/cache"
	"github.com/niklod/highload-social-network/internal/queue/feed/producer"
)

const (
	feedLength = 1000
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
	repo     repository
	cache    cache.Cache
	producer *producer.FeedProducer
}

func NewService(repo repository, cache cache.Cache, producer *producer.FeedProducer) *Service {
	return &Service{
		repo:     repo,
		cache:    cache,
		producer: producer,
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
		if len(userFeed) > feedLength {
			userFeed = userFeed[len(userFeed)-feedLength:]
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

func (s *Service) Add(post *Post) error {
	if post == nil {
		return errNilPost
	}

	authorId := post.Author.ID

	if authorId <= 0 {
		return errIdLessThanZero
	}
	if post.Body == "" {
		return errEmptyPostBody
	}

	err := s.repo.Add(post, authorId)
	if err != nil {
		return fmt.Errorf("post.Service: %v", err)
	}

	post.CreatedAt = time.Now().UTC()

	msg, err := post.AsByteJSON()
	if err != nil {
		return fmt.Errorf("post.Service - can't marshal post to []byte: %v", err)
	}

	err = s.producer.SendFeedUpdateMessage(msg)
	if err != nil {
		return fmt.Errorf("post.Service - can't send message to queue: %v", err)
	}

	return nil
}
