package repositories

import "catecard/pkg/domain/entities"

type UserRepository interface {
	GetAll() ([]*entities.User, error)
	GetUser(username string, password string) *entities.User
	SaveUser(user *entities.User) (*entities.User, error)
}
type ProductRepository interface {
	Add(product *entities.Product) error
	GetAll() ([]*entities.Product, error)
}
