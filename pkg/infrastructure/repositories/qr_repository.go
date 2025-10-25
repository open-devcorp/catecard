package repositories

import (
	"catecard/pkg/domain/entities"
	"database/sql"
	"log"
)

type qrRepository struct {
	log *log.Logger
	db  *sql.DB
}

func (q *qrRepository) Add(qr *entities.Qr) (int, error) {
	query := `INSERT INTO qrs (group_id, forum, count) VALUES (?, ?, ?)`
	result, err := q.db.Exec(query, qr.GroupId, qr.Forum, qr.Count)
	if err != nil {
		q.log.Printf("Error inserting QR: %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		q.log.Printf("Error getting last insert ID: %v", err)
		return 0, err
	}

	qr.ID = int(id)
	return qr.ID, nil
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

// ClaimAtomic realiza la actualización de cupo de forma segura:
// Solo decrementa forum e incrementa count si forum > 0.
// Protege contra carreras cuando dos escaneos llegan "al mismo tiempo".
func (q *qrRepository) ClaimAtomic(id int) (bool, error) {
	res, err := q.db.Exec(`
		UPDATE qrs
		SET forum = forum - 1,
		    count = count + 1
		WHERE id = ? AND forum > 0
	`, id)
	if err != nil {
		q.log.Printf("Error in ClaimAtomic UPDATE: %v", err)
		return false, err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		q.log.Printf("Error getting RowsAffected in ClaimAtomic: %v", err)
		return false, err
	}
	// aff == 0 => no había cupo (forum <= 0)
	return aff > 0, nil
}

// Update: mantenla si te sirve para cambios administrativos (no para consumir cupos).
// Ya NO valida forum aquí; la validación de cupo se hace en ClaimAtomic.
func (q *qrRepository) Update(qr *entities.Qr) error {
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
