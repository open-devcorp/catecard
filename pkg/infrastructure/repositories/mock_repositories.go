package repositories

import (
	"catecard/pkg/domain/entities"
)

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
