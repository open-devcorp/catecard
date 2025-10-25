package repositories

import (
	"catecard/pkg/domain/entities"
	"database/sql"
	"log"
)

type scanCatechumenRepository struct {
	log *log.Logger
	db  *sql.DB
}

func (s *scanCatechumenRepository) GetAll() ([]ScanAndCatechume, error) {
	query := `
		SELECT
			c.id, c.full_name, c.age, c.group_id,
			u.id, u.username, u.email, u.password, u.role
		FROM scan_catechumens sc
		JOIN catechumens c ON sc.catechumen_id = c.id
		JOIN users u       ON sc.scan_id = u.id
	`

	rows, err := s.db.Query(query)
	if err != nil {
		s.log.Printf("Error ejecutando la consulta: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []ScanAndCatechume

	for rows.Next() {
		var catechumen entities.Catechumen
		var user entities.User

		// OJO: age es VARCHAR(3) en la BD → usa string en la entidad
		if err := rows.Scan(
			&catechumen.ID,
			&catechumen.FullName,
			&catechumen.Age, // string si en BD es VARCHAR(3)
			&catechumen.GroupId,
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.Role,
		); err != nil {
			s.log.Printf("Error escaneando fila: %v", err)
			return nil, err
		}

		results = append(results, ScanAndCatechume{
			User:       &user,
			Catechumen: &catechumen,
		})
	}

	if err := rows.Err(); err != nil {
		s.log.Printf("Error al recorrer filas: %v", err)
		return nil, err
	}

	return results, nil
}

func NewScanCatechumenRepository(logger *log.Logger, db *sql.DB) ScanCatechumenRepository {
	return &scanCatechumenRepository{logger, db}
}

func (s *scanCatechumenRepository) Add(scanCatechumen *entities.ScanCatechumen) error {
	query := `INSERT INTO scan_catechumens (catechumen_id, scan_id, created_at) VALUES (?, ?, ?)`
	_, err := s.db.Exec(query, scanCatechumen.CatechumenID, scanCatechumen.ScanID, scanCatechumen.CreatedAt)
	if err != nil {
		s.log.Printf("Error adding scan catechumen: %v", err)
		return err
	}
	return nil
}
