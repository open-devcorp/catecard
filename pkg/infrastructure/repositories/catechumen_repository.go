package repositories

import (
	"catecard/pkg/domain/entities"
	"database/sql"
	"log"
)

type catechumenRepository struct {
	log *log.Logger
	db  *sql.DB
}

// Add implements CatechumenRepository.
func (c *catechumenRepository) Add(catechumen *entities.Catechumen) (int, error) {
	query := `INSERT INTO catechumens(full_name, age, group_id) VALUES(?,?,?)`
	result, err := c.db.Exec(query, catechumen.FullName, catechumen.Age, catechumen.GroupId)

	if err != nil {
		c.log.Printf("Error inserting catechumen: %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.log.Printf("Error getting last insert ID: %v", err)
		return 0, err
	}
	catechumen.ID = int(id)
	return catechumen.ID, nil
}

// GetAll implements CatechumenRepository.
func (c *catechumenRepository) GetAll() ([]*entities.Catechumen, error) {
	query := `SELECT id, full_name, age, group_id FROM catechumens`
	rows, err := c.db.Query(query)

	if err != nil {
		c.log.Printf("Error querying catechumens: %v", err)
		return nil, err
	}

	defer rows.Close()
	var catechumens []*entities.Catechumen
	for rows.Next() {
		catechumen := &entities.Catechumen{}
		if err := rows.Scan(&catechumen.ID, &catechumen.FullName, &catechumen.Age, &catechumen.GroupId); err != nil {
			c.log.Printf("Error scanning catechumen: %v", err)
			return nil, err
		}
		catechumens = append(catechumens, catechumen)
	}

	return catechumens, nil
}

// GetById implements CatechumenRepository.
func (c *catechumenRepository) GetById(id int) (*entities.Catechumen, error) {
	panic("unimplemented")
}

// Update implements CatechumenRepository.
func (c *catechumenRepository) Update(catechumen *entities.Catechumen) (*entities.Catechumen, error) {
	panic("unimplemented")
}

func NewCatechumenRepository(logger *log.Logger, db *sql.DB) CatechumenRepository {
	return &catechumenRepository{logger, db}
}
