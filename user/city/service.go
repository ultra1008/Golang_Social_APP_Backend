package city

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(city City) (*City, error) {
	return s.repo.Create(city.Name)
}

func (s *Service) List() ([]City, error) {
	return s.repo.List()
}
