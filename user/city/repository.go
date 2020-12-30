package city

type repository interface {
	Create(city string) (*City, error)
	List() ([]City, error)
}
