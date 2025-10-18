package repositories

import (
	"catecard/pkg/domain/entities"
	"database/sql"
	"log"
)

type groupRepository struct {
	log *log.Logger
	db  *sql.DB
}

func NewGroupRepository(logger *log.Logger, db *sql.DB) GroupRepository {
	return &groupRepository{logger, db}
}

// Add implements GroupRepository.
func (g *groupRepository) Add(group *entities.Group) error {

	query := `INSERT INTO groups(name, catechist_id) VALUES (?,?)`
	result, err := g.db.Exec(query, group.Name, group.CatechistId)
	if err != nil {
		g.log.Printf("Error inserting group: %v", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		g.log.Printf("Error getting last insert ID: %v", err)
		return err
	}
	group.ID = int(id)
	return nil
}

// Edit implements GroupRepository.
func (g *groupRepository) Edit(group *entities.Group) error {
	panic("unimplemented")
}

// GetAll implements GroupRepository.
func (g *groupRepository) GetAll() ([]*entities.Group, error) {
	panic("unimplemented")
}
