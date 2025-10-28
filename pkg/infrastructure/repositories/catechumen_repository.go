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
	// Adaptado a esquema qr_codes: buscar catechumen por join qr_codes -> catechumens
	query := `SELECT c.id, c.full_name, c.age, c.group_id
			  FROM catechumens c
			  JOIN qr_codes q ON q.catechumen_id = c.id
			 WHERE q.id = ?`
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

// DeleteById implements CatechumenRepository.
func (c *catechumenRepository) DeleteById(id int) error {
	query := `DELETE FROM catechumens WHERE id = ?`
	_, err := c.db.Exec(query, id)

	if err != nil {
		c.log.Printf("Error deleting catechumen by ID: %v", err)
		return err
	}

	return nil
}

// Add implements CatechumenRepository.
func (c *catechumenRepository) Add(catechumen *entities.Catechumen) (int, error) {
	// Esquema actual: catechumens ya no lleva qr_id; QR se crea en qr_codes con catechumen_id
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
	// Trae catecúmeno y, si existe, el qr_id derivado desde qr_codes
	query := `SELECT c.id, c.full_name, c.age, c.group_id, COALESCE(q.id, 0) AS qr_id
			  FROM catechumens c
			  LEFT JOIN qr_codes q ON q.catechumen_id = c.id
			 WHERE c.id = ?`
	row := c.db.QueryRow(query, id)

	catechumen := &entities.Catechumen{}
	if err := row.Scan(&catechumen.ID, &catechumen.FullName, &catechumen.Age, &catechumen.GroupId, &catechumen.QrId); err != nil {
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
