package usecases

import (
	"catecard/pkg/domain/entities"
	"catecard/pkg/infrastructure/repositories"
	"fmt"
	"log"
)

type CatechumenUseCase interface {
	Add(catechumen *entities.Catechumen) (*entities.Catechumen, *entities.Qr, error)
	Update(catechumen *entities.Catechumen) (*entities.Catechumen, error)
	GetAll() ([]*entities.Catechumen, error)
	GetById(id int) (*entities.Catechumen, error)
}

type catechumenUseCase struct {
	logger         *log.Logger
	catechumenRepo repositories.CatechumenRepository
	qrRepo         repositories.QrRepository
}

func NewCatechumenUsecase(logger *log.Logger, catechumenRepo repositories.CatechumenRepository, qrRepo repositories.QrRepository) CatechumenUseCase {
	return &catechumenUseCase{logger, catechumenRepo, qrRepo}
}

func (c *catechumenUseCase) Add(catechumen *entities.Catechumen) (*entities.Catechumen, *entities.Qr, error) {
	if catechumen.FullName == "" {
		return nil, nil, fmt.Errorf("Group name cannot be empty")
	}
	if catechumen.Age == "" {
		return nil, nil, fmt.Errorf("Catechumen age cannot be empty")
	}
	if catechumen.GroupId == 0 {
		return nil, nil, fmt.Errorf("Group ID cannot be 0")
	}

	cateId, err := c.catechumenRepo.Add(catechumen)
	if err != nil {
		c.logger.Printf("Error adding catechumen: %v", err)
		return nil, nil, err
	}

	qr := entities.NewQr(3, cateId)

	qrErr := c.qrRepo.Add(qr)
	if qrErr != nil {
		c.logger.Printf("Error adding QR: %v", qrErr)
		return nil, nil, qrErr
	}

	return catechumen, qr, nil
}

func (c *catechumenUseCase) Update(catechumen *entities.Catechumen) (*entities.Catechumen, error) {
	if catechumen.FullName == "" {
		return nil, fmt.Errorf("Group name cannot be empty")
	}
	if catechumen.Age == "" {
		return nil, fmt.Errorf("Catechumen age cannot be empty")
	}
	if catechumen.GroupId == 0 {
		return nil, fmt.Errorf("Group ID cannot be 0")
	}

	updatedCatechumen, err := c.catechumenRepo.Update(catechumen)
	if err != nil {
		c.logger.Printf("Error updating catechumen: %v", err)
		return nil, err
	}

	return updatedCatechumen, nil
}

func (c *catechumenUseCase) GetAll() ([]*entities.Catechumen, error) {
	if c.catechumenRepo == nil {
		return nil, fmt.Errorf("Catechumen repository is not initialized")
	}
	return c.catechumenRepo.GetAll()
}

func (c *catechumenUseCase) GetById(id int) (*entities.Catechumen, error) {
	if id == 0 {
		return nil, fmt.Errorf("Invalid catechumen ID")
	}
	catechumen, err := c.catechumenRepo.GetById(id)
	if err != nil {
		c.logger.Printf("Error getting catechumen by ID: %v", err)
		return nil, err
	}
	return catechumen, nil
}
