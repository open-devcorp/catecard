package usecases

import (
	"catecard/pkg/domain/entities"
	"catecard/pkg/infrastructure/repositories"
	"fmt"
	"log"
)

type GroupUseCase interface {
	Add(group *entities.Group) error
	GetAll(User *entities.User) ([]*entities.Group, error)
	GetById(User *entities.User, id int) (*entities.Group, error)
	DeleteById(User *entities.User, id int) error
	Update(User *entities.User, group *entities.Group) (*entities.Group, error)
}

type groupUseCase struct {
	logger    *log.Logger
	groupRepo repositories.GroupRepository
}

func NewGroupUsecase(logger *log.Logger, r repositories.GroupRepository) GroupUseCase {
	return &groupUseCase{logger, r}
}

// Update implements GroupUseCase.
func (g *groupUseCase) Update(User *entities.User, group *entities.Group) (*entities.Group, error) {
	if User == nil {
		return nil, fmt.Errorf("unauthorized: no user in session")
	}
	if User.Role != entities.ADMIN {
		return nil, fmt.Errorf("forbidden: user does not have admin role")
	}

	if group.Name == "" {
		return nil, fmt.Errorf("Group name cannot be empty")
	}
	if group.CatechistId == 0 {
		return nil, fmt.Errorf("Catechist ID cannot be 0")
	}

	return g.groupRepo.Update(group)
}

// DeleteById implements GroupUseCase.
func (g *groupUseCase) DeleteById(User *entities.User, id int) error {
	if User == nil {
		return fmt.Errorf("unauthorized: no user in session")
	}
	if User.Role != entities.ADMIN {
		return fmt.Errorf("forbidden: user does not have admin role")
	}

	return g.groupRepo.DeleteById(id)
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
