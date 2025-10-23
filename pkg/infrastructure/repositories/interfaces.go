package repositories

import "catecard/pkg/domain/entities"

type UserRepository interface {
	GetAll() ([]*entities.User, error)
	GetUser(username string, password string) *entities.User
	SaveUser(user *entities.User) (*entities.User, error)
	DeleteUserById(id int) error
	GetById(id int) (*entities.User, error)
}

type GroupRepository interface {
	Add(group *entities.Group) error
	GetAll() ([]*entities.Group, error)
	GetById(id int) (*entities.Group, error)
	DeleteById(id int) error
	Update(group *entities.Group) (*entities.Group, error)
	GetByCatechistsId(catechistId int) (int, error)
}

type QrRepository interface {
	Add(qr *entities.Qr) error
	GetAll() ([]*entities.Qr, error)
	GetById(id int) (*entities.Qr, error)
	DeleteById(id int) error
	Update(qr *entities.Qr) error
	ClaimAtomic(id int) (bool, error)
}

type CatechumenRepository interface {
	Add(catechumen *entities.Catechumen) (int, error)
	Update(catechumen *entities.Catechumen) (*entities.Catechumen, error)
	GetAll() ([]*entities.Catechumen, error)
	GetById(id int) (*entities.Catechumen, error)
}
