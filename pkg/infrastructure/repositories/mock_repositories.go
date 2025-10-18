package repositories

import (
	"catecard/pkg/domain/entities"
)

type MockProductRepository struct {
	Products []*entities.Product
}

func NewMockProductRepository() *MockProductRepository {

	return &MockProductRepository{
		Products: make([]*entities.Product, 0),
	}
}

// SIMULA LA DB
func (m *MockProductRepository) Add(product *entities.Product) error {

	m.Products = append(m.Products, product)
	return nil
}

func (m *MockProductRepository) GetAll() ([]*entities.Product, error) {
	product := entities.FakeProduct()
	m.Products = append(m.Products, &product)

	return m.Products, nil
}

type MockUserRepository struct {
	Users []*entities.User
}

// DeleteUserById implements UserRepository.
func (m *MockUserRepository) DeleteUserById(id int) error {
	panic("unimplemented")
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{Users: make([]*entities.User, 0)}
}

func (m *MockUserRepository) GetAll() ([]*entities.User, error) {
	if len(m.Users) == 0 {
		u := entities.FakeUser()
		m.Users = append(m.Users, &u)
	}
	return m.Users, nil
}

func (m *MockUserRepository) GetUser(username string, password string) *entities.User {
	for _, u := range m.Users {
		if u.Username == username && u.Password == password {
			return u
		}
	}
	return nil
}

func (m *MockUserRepository) SaveUser(user *entities.User) (*entities.User, error) {
	user.ID = len(m.Users) + 1
	m.Users = append(m.Users, user)

	return user, nil
}
