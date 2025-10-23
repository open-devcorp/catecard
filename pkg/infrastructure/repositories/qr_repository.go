package repositories

import (
	"catecard/pkg/domain/entities"
	"database/sql"
	"errors"
	"log"
)

type qrRepository struct {
	log *log.Logger
	db  *sql.DB
}

func (q *qrRepository) Add(qr *entities.Qr) error {
	query := `INSERT INTO qrs (group_id, forum, count) VALUES (?, ?, ?)`
	result, err := q.db.Exec(query, qr.GroupId, qr.Forum, qr.Count)
	if err != nil {
		q.log.Printf("Error inserting QR: %v", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		q.log.Printf("Error getting last insert ID: %v", err)
		return err
	}

	qr.ID = int(id)
	return nil
}

func (q *qrRepository) DeleteById(id int) error {
	query := `DELETE FROM qrs WHERE id = ?`
	_, err := q.db.Exec(query, id)
	if err != nil {
		q.log.Printf("Error deleting QR by ID: %v", err)
		return err
	}
	return nil
}

func (q *qrRepository) GetAll() ([]*entities.Qr, error) {
	query := `SELECT id, group_id, forum, count FROM qrs`
	rows, err := q.db.Query(query)
	if err != nil {
		q.log.Printf("Error getting all QRs: %v", err)
		return nil, err
	}
	defer rows.Close()

	var qrs []*entities.Qr
	for rows.Next() {
		var qr entities.Qr
		if err := rows.Scan(&qr.ID, &qr.GroupId, &qr.Forum, &qr.Count); err != nil {
			q.log.Printf("Error scanning QR: %v", err)
			return nil, err
		}
		qrs = append(qrs, &qr)
	}
	return qrs, nil
}

// GetById obtiene un QR específico por su ID.
func (q *qrRepository) GetById(id int) (*entities.Qr, error) {
	query := `SELECT id, group_id, forum, count FROM qrs WHERE id = ?`
	row := q.db.QueryRow(query, id)

	var qr entities.Qr
	if err := row.Scan(&qr.ID, &qr.GroupId, &qr.Forum, &qr.Count); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		q.log.Printf("Error getting QR by ID: %v", err)
		return nil, err
	}
	return &qr, nil
}

func (q *qrRepository) Update(qr *entities.Qr) error {
	if qr.Forum <= 0 {
		return errors.New("Se alcanzó el límite máximo de participantes para este QR")
	}

	query := `UPDATE qrs SET group_id = ?, forum = ?, count = ? WHERE id = ?`
	_, err := q.db.Exec(query, qr.GroupId, qr.Forum, qr.Count, qr.ID)
	if err != nil {
		q.log.Printf("Error updating QR: %v", err)
		return err
	}

	return nil
}

func NewQrRepository(logger *log.Logger, db *sql.DB) QrRepository {
	return &qrRepository{logger, db}
}
