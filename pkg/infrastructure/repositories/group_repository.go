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
	query := `SELECT id, name, catechist_id FROM groups`
	rows, err := g.db.Query(query)
	if err != nil {
		g.log.Printf("Error querying groups: %v", err)
		return nil, err
	}
	defer rows.Close()

	var groups []*entities.Group
	for rows.Next() {
		group := &entities.Group{}
		err := rows.Scan(&group.ID, &group.Name, &group.CatechistId)
		if err != nil {
			g.log.Printf("Error scanning group: %v", err)
			return nil, err
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		g.log.Printf("Error iterating over group rows: %v", err)
		return nil, err
	}

	return groups, nil

}

func (g *groupRepository) GetById(id int) (*entities.Group, error) {
	query := `SELECT id, name, catechist_id FROM groups WHERE id = ?`
	row := g.db.QueryRow(query, id)

	group := &entities.Group{}
	err := row.Scan(&group.ID, &group.Name, &group.CatechistId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No group found with the given ID
		}
		g.log.Printf("Error getting group by ID: %v", err)
		return nil, err
	}

	return group, nil
}
