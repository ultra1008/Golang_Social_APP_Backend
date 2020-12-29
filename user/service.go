package user

type Service struct {
	db repository
}

func NewService(repo repository) *Service {
	return &Service{
		db: repo,
	}
}
