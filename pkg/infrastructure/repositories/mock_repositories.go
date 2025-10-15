package repositories

import "goapp/pkg/domain/entities"

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
