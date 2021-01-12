package user

import (
	"fmt"
	"log"

	"github.com/niklod/highload-social-network/user/city"
	"github.com/niklod/highload-social-network/user/interest"
	"golang.org/x/crypto/bcrypt"
)

type repository interface {
	Create(user *User) (*User, error)
	List() ([]User, error)
	GetByID(id int) (*User, error)
	GetByLogin(login string) (*User, error)
	AddFriend(userId int, friendId int) error
	DeleteFriend(userId int, friendId int) error
	Friends(userId int) ([]User, error)
}

type Service struct {
	userRepo        repository
	cityService     *city.Service
	interestService *interest.Service
}

func NewService(repo repository, citySvc *city.Service, interestSvc *interest.Service) *Service {
	return &Service{
		userRepo:        repo,
		cityService:     citySvc,
		interestService: interestSvc,
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

	fmt.Printf("%+v", user.Interests)

	hash, err := s.CreatePassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("generating hash from password: %v", err)
	}

	user.Password = hash

	updatedUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	for _, i := range user.Interests {
		err := s.interestService.Create(&i)
		if err != nil {
			log.Printf("creating user interest %v", err)
			continue
		}

		err = s.interestService.AddInterestToUser(updatedUser.ID, i.ID)
		if err != nil {
			log.Printf("addint interest to user %v", err)
			continue
		}
	}

	fmt.Printf("%+v", updatedUser)

	return updatedUser, nil
}

func (s *Service) Users() ([]User, error) {
	return s.userRepo.List()
}

func (s *Service) CreatePassword(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generating hash from password: %v", err)
	}

	return string(hash), nil
}

func (s *Service) CheckPasswordsEquality(pass, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
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

func (s *Service) GetUserByLogin(userLogin string) (*User, error) {
	return s.userRepo.GetByLogin(userLogin)
}

func (s *Service) AddFriend(userId, friendId int) error {
	return s.userRepo.AddFriend(userId, friendId)
}

func (s *Service) DeleteFriend(userId, friendId int) error {
	return s.userRepo.DeleteFriend(userId, friendId)
}

func (s *Service) Friends(userId int) ([]User, error) {
	return s.userRepo.Friends(userId)
}

func (s *Service) IsUsersAreFriends(user, userToCheck *User) bool {
	if user == nil {
		return false
	}

	for _, f := range user.Friends {
		if userToCheck.ID == f.ID {
			return true
		}
	}
	return false
}

func (s *Service) Interests(userId int) ([]interest.Interest, error) {
	return s.interestService.Interests()
}

func (s *Service) AddInterest(userId, interestId int) error {
	return s.interestService.AddInterestToUser(userId, interestId)
}
