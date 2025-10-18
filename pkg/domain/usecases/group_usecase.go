package usecases

import (
	"catecard/pkg/domain/entities"
	"catecard/pkg/infrastructure/repositories"
	"fmt"
	"log"
)

type GroupUseCase interface {
	Add(group *entities.Group) error
	Edit(group *entities.Group) error
	GetAll(User *entities.User) ([]*entities.Group, error)
	GetById(User *entities.User, id int) (*entities.Group, error)
}

func NewGroupUsecase(logger *log.Logger, r repositories.GroupRepository) GroupUseCase {
	return &groupUseCase{logger, r}
}

type groupUseCase struct {
	logger    *log.Logger
	groupRepo repositories.GroupRepository
}

// GetById implements GroupUseCase.
func (g *groupUseCase) GetById(User *entities.User, id int) (*entities.Group, error) {

	if User == nil {
		return nil, fmt.Errorf("unauthorized: no user in session")
	}
	if g.groupRepo == nil {
		return nil, fmt.Errorf("Group repository is not initialized")
	}
	return g.groupRepo.GetById(id)

}

// Edit implements GroupUseCase.
func (g *groupUseCase) Edit(group *entities.Group) error {
	panic("unimplemented")
}

// GetAll implements GroupUseCase.
func (g *groupUseCase) GetAll(User *entities.User) ([]*entities.Group, error) {
	if g.groupRepo == nil {
		return nil, fmt.Errorf("Group repository is not initialized")
	}
	return g.groupRepo.GetAll()
}

// Add implements GroupUseCase.
func (g *groupUseCase) Add(group *entities.Group) error {

	if group.Name == "" {
		return fmt.Errorf("Group name cannot be empty")
	}
	if group.CatechistId == 0 {
		return fmt.Errorf("Catechist ID cannot be 0")
	}

	err := g.groupRepo.Add(group)
	if err != nil {
		g.logger.Printf("Error adding group: %v", err)
		return err
	}
	return nil
}
