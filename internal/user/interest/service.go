package interest

type repository interface {
	CreateIfNotExists(i *Interest) error
	List() ([]Interest, error)
	InterestsByUserId(id int) ([]Interest, error)
	AddInterestToUser(userId, interestId int) error
}

type Service struct {
	InterestRepo repository
}

func NewService(repo repository) *Service {
	return &Service{
		InterestRepo: repo,
	}
}

func (s *Service) Interests() ([]Interest, error) {
	return s.InterestRepo.List()
}

func (s *Service) InterestsByUserId(id int) ([]Interest, error) {
	return s.InterestRepo.InterestsByUserId(id)
}

func (s *Service) Create(i *Interest) error {
	return s.InterestRepo.CreateIfNotExists(i)
}

func (s *Service) AddInterestToUser(userId, interestId int) error {
	return s.InterestRepo.AddInterestToUser(userId, interestId)
}
