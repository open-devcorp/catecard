package repositories

import (
	"catecard/pkg/domain/entities"
	"database/sql"
	"fmt"
	"log"
)

type qrRepository struct {
	log *log.Logger
	db  *sql.DB
}

func (q *qrRepository) Add(qr *entities.Qr) (int, error) {
	// Adaptado a esquema qr_codes: requiere catechumen_id y totales
	// Usamos qr.Catechumen.ID como foreign key y mapear:
	// total_allowed = qr.Forum (capacidad inicial)
	// used_scans   = qr.Count  (usos iniciales, usualmente 0)
	if qr.Catechumen == nil || qr.Catechumen.ID == 0 {
		q.log.Printf("Add QR: missing catechumen ID")
		return 0, fmt.Errorf("missing catechumen ID for QR creation")
	}
	query := `INSERT INTO qr_codes (catechumen_id, total_allowed, used_scans) VALUES (?, ?, ?)`
	result, err := q.db.Exec(query, qr.Catechumen.ID, qr.Forum, qr.Count)
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
	query := `DELETE FROM qr_codes WHERE id = ?`
	_, err := q.db.Exec(query, id)
	if err != nil {
		q.log.Printf("Error deleting QR by ID: %v", err)
		return err
	}
	return nil
}

func (q *qrRepository) GetAll() ([]*entities.Qr, error) {
	// Mapear desde qr_codes: forum = total_allowed - used_scans, count = used_scans
	query := `SELECT id, catechumen_id, total_allowed, used_scans FROM qr_codes`
	rows, err := q.db.Query(query)
	if err != nil {
		q.log.Printf("Error getting all QRs: %v", err)
		return nil, err
	}
	defer rows.Close()

	var qrs []*entities.Qr
	for rows.Next() {
		var (
			id     int
			cateId int
			total  int
			used   int
		)
		if err := rows.Scan(&id, &cateId, &total, &used); err != nil {
			q.log.Printf("Error scanning QR: %v", err)
			return nil, err
		}
		qr := &entities.Qr{ID: id, Forum: total - used, Count: used}
		if cateId != 0 {
			qr.Catechumen = &entities.Catechumen{ID: cateId}
		}
		qrs = append(qrs, qr)
	}
	return qrs, nil
}

// GetById obtiene un QR específico por su ID.
func (q *qrRepository) GetById(id int) (*entities.Qr, error) {
	// Mapear desde qr_codes: forum = total_allowed - used_scans, count = used_scans
	query := `SELECT id, catechumen_id, total_allowed, used_scans FROM qr_codes WHERE id = ?`
	row := q.db.QueryRow(query, id)

	var (
		qr     entities.Qr
		cateId int
		total  int
		used   int
	)
	if err := row.Scan(&qr.ID, &cateId, &total, &used); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		q.log.Printf("Error getting QR by ID: %v", err)
		return nil, err
	}
	qr.Forum = total - used
	qr.Count = used
	if cateId != 0 {
		qr.Catechumen = &entities.Catechumen{ID: cateId}
	}
	return &qr, nil
}

// ClaimAtomic realiza la actualización de cupo de forma segura:
// Solo decrementa forum e incrementa count si forum > 0.
// Protege contra carreras cuando dos escaneos llegan "al mismo tiempo".
func (q *qrRepository) ClaimAtomic(id int) (bool, error) {
	// Incrementa used_scans si aún no alcanzó total_allowed
	res, err := q.db.Exec(`
		UPDATE qr_codes
		SET used_scans = used_scans + 1
		WHERE id = ? AND used_scans < total_allowed
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
	// En qr_codes no hay group_id ni forum directo; actualizamos derived fields:
	// total_allowed = qr.Forum + qr.Count (si se modifica cupo) y used_scans = qr.Count
	// Para no romper usos existentes (claim), bastará con reflejar used_scans = qr.Count.
	query := `UPDATE qr_codes SET used_scans = ? WHERE id = ?`
	_, err := q.db.Exec(query, qr.Count, qr.ID)
	if err != nil {
		q.log.Printf("Error updating QR: %v", err)
		return err
	}
	return nil
}

func NewQrRepository(logger *log.Logger, db *sql.DB) QrRepository {
	return &qrRepository{logger, db}
}
