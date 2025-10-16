package usecases

import (
	"catecard/pkg/domain/entities"
	"catecard/pkg/infrastructure/repositories"
	"fmt"
	"log"
)

type ProductUseCase interface {
	Add(product *entities.Product) error
	GetAll() ([]*entities.Product, error)
}

type productUseCase struct {
	logger      *log.Logger
	productRepo repositories.ProductRepository
}

// Add implements ProductUseCase.
func (p *productUseCase) Add(product *entities.Product) error {

	//1
	if product.Name == "" {
		return fmt.Errorf("Product name cannot be empty")
	}
	if product.Price == 0.0 {
		return fmt.Errorf("Product price cannot be 0")
	}

	//2
	err := p.productRepo.Add(product)
	if err != nil {
		p.logger.Printf("Error adding product: %v", err)
		return err
	}
	return nil

}

// GetAll implements ProductUseCase.
func (p *productUseCase) GetAll() ([]*entities.Product, error) {

	products, err := p.productRepo.GetAll()

	if err != nil {

		p.logger.Printf("Error getting all products: %v", err)
		return nil, err
	}
	return products, nil

}

func NewProductUsecase(logger *log.Logger, r repositories.ProductRepository) ProductUseCase {
	return &productUseCase{logger, r}

}
