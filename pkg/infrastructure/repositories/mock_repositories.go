package repositories

import (
	"catecard/pkg/domain/entities"
)

// MockUserRepository is a mock implementation of UserRepository for testing purposes.
type MockUserRepository struct {
	Users []*entities.User
}

// NewMockUserRepository creates a new instance of MockUserRepository.
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{Users: make([]*entities.User, 0)}
}

func (m *MockUserRepository) DeleteUserById(id int) error {
	panic("unimplemented")
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

// MockCatechumenRepository is a mock implementation of CatechumenRepository for testing purposes.
type MockCatechumenRepository struct {
	Catechumens []*entities.Catechumen
}

// NewMockCatechumenRepository creates a new instance of MockCatechumenRepository.
func NewMockCatechumenRepository() *MockCatechumenRepository {
	return &MockCatechumenRepository{Catechumens: make([]*entities.Catechumen, 0)}
}

func (m *MockCatechumenRepository) Add(catechumen *entities.Catechumen) (int, error) {
	catechumen.ID = len(m.Catechumens) + 1
	m.Catechumens = append(m.Catechumens, catechumen)
	return catechumen.ID, nil
}

func (m *MockCatechumenRepository) Update(catechumen *entities.Catechumen) (*entities.Catechumen, error) {
	for i, c := range m.Catechumens {
		if c.ID == catechumen.ID {
			m.Catechumens[i] = catechumen
			return catechumen, nil
		}
	}
	return nil, nil
}

func (m *MockCatechumenRepository) GetAll() ([]*entities.Catechumen, error) {
	if len(m.Catechumens) == 0 {
		c := entities.FakeCatechumen()
		m.Catechumens = append(m.Catechumens, &c)
	}
	return m.Catechumens, nil
}

func (m *MockCatechumenRepository) GetById(id int) (*entities.Catechumen, error) {
	for _, c := range m.Catechumens {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, nil
}

// MockQrRepository is a mock implementation of QrRepository for testing purposes.
type MockQrRepository struct {
	Qrs []*entities.Qr
}

// NewMockQrRepository creates a new instance of MockQrRepository.
func NewMockQrRepository() *MockQrRepository {
	return &MockQrRepository{Qrs: make([]*entities.Qr, 0)}
}

func (m *MockQrRepository) Add(qr *entities.Qr) error {
	qr.ID = len(m.Qrs) + 1
	m.Qrs = append(m.Qrs, qr)
	return nil
}

func (m *MockQrRepository) GetAll() ([]*entities.Qr, error) {
	if len(m.Qrs) == 0 {
		q := entities.FakeQr()
		m.Qrs = append(m.Qrs, &q)
	}
	return m.Qrs, nil
}

func (m *MockQrRepository) GetById(id int) (*entities.Qr, error) {
	for _, q := range m.Qrs {
		if q.ID == id {
			return q, nil
		}
	}
	return nil, nil
}

func (m *MockQrRepository) DeleteById(id int) error {
	panic("unimplemented")
}

func (m *MockQrRepository) Update(qr *entities.Qr) error {
	panic("unimplemented")
}
