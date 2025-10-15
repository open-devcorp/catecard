package repositories

import "goapp/pkg/domain/entities"

type ProductRepository interface {
	Add(product *entities.Product) error
	GetAll() ([]*entities.Product, error)
}
