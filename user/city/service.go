package city

type Serivce struct {
	repo repository
}

func NewService(repo repository) *Serivce {
	return &Serivce{
		repo: repo,
	}
}

func (s *Serivce) Create(city City) (*City, error) {
	return s.repo.Create(city.Name)
}
