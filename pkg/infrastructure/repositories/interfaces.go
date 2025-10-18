package repositories

import "catecard/pkg/domain/entities"

type UserRepository interface {
	GetAll() ([]*entities.User, error)
	GetUser(username string, password string) *entities.User
	SaveUser(user *entities.User) (*entities.User, error)
	DeleteUserById(id int) error
}
type ProductRepository interface {
	Add(product *entities.Product) error
	GetAll() ([]*entities.Product, error)
}

type GroupRepository interface {
	Add(group *entities.Group) error
	Edit(group *entities.Group) error
	GetAll() ([]*entities.Group, error)
	GetById(id int) (*entities.Group, error)
}
