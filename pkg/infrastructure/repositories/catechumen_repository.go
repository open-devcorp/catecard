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

// GetByQrId implements CatechumenRepository.
func (c *catechumenRepository) GetByQrId(qrId int) (*entities.Catechumen, error) {
	query := `SELECT id, full_name, age, group_id FROM catechumens WHERE qr_id = ?`
	row := c.db.QueryRow(query, qrId)

	catechumen := &entities.Catechumen{}
	if err := row.Scan(&catechumen.ID, &catechumen.FullName, &catechumen.Age, &catechumen.GroupId); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		c.log.Printf("Error scanning catechumen by QR ID: %v", err)
		return nil, err
	}

	return catechumen, nil
}

// Add implements CatechumenRepository.
func (c *catechumenRepository) Add(catechumen *entities.Catechumen) (int, error) {
	query := `INSERT INTO catechumens(full_name, age, group_id, qr_id) VALUES(?,?,?,?)`
	result, err := c.db.Exec(query, catechumen.FullName, catechumen.Age, catechumen.GroupId, catechumen.QrId)

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
	query := `SELECT id, full_name, age, group_id FROM catechumens WHERE id = ?`
	row := c.db.QueryRow(query, id)

	catechumen := &entities.Catechumen{}
	if err := row.Scan(&catechumen.ID, &catechumen.FullName, &catechumen.Age, &catechumen.GroupId); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		c.log.Printf("Error scanning catechumen by ID: %v", err)
		return nil, err
	}

	return catechumen, nil
}

// Update implements CatechumenRepository.
func (c *catechumenRepository) Update(catechumen *entities.Catechumen) (*entities.Catechumen, error) {
	query := `UPDATE catechumens SET full_name = ?, age = ?, group_id = ? WHERE id = ?`
	_, err := c.db.Exec(query, catechumen.FullName, catechumen.Age, catechumen.GroupId, catechumen.ID)

	if err != nil {
		c.log.Printf("Error updating catechumen: %v", err)
		return nil, err
	}

	return catechumen, nil
}

func NewCatechumenRepository(logger *log.Logger, db *sql.DB) CatechumenRepository {
	return &catechumenRepository{logger, db}
}
