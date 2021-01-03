package user

import (
	"fmt"

	"github.com/niklod/highload-social-network/user/city"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo    repository
	cityService *city.Serivce
}

func NewService(repo repository, citySvc *city.Serivce) *Service {
	return &Service{
		userRepo:    repo,
		cityService: citySvc,
	}
}

func (s *Service) Create(user *User) (*User, error) {
	ok, err := s.CheckUserExist(user.Login)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, fmt.Errorf("user already exist")
	}

	city, err := s.cityService.Create(user.City)
	if err != nil {
		return nil, err
	}

	user.City = *city

	hash, err := s.CreatePassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("generating hash from password: %v", err)
	}

	user.Password = hash

	updatedUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *Service) CreatePassword(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generating hash from password: %v", err)
	}

	return string(hash), nil
}

func (s *Service) CheckUserExist(userLogin string) (bool, error) {
	user, err := s.userRepo.GetByLogin(userLogin)
	if err != nil {
		return false, err
	}

	if user != nil {
		return true, nil
	}

	return false, nil
}
