package handlers

import (
	"catecard/pkg/domain/usecases"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type QrHandler interface {
	AddQr(w http.ResponseWriter, r *http.Request)
	GetAllQrs(w http.ResponseWriter, r *http.Request)
	GetQrById(id int, w http.ResponseWriter, r *http.Request)
	ClaimQr(id int, w http.ResponseWriter, r *http.Request)
}

type qrHandler struct {
	log      *log.Logger
	uc       usecases.QrUseCase
	tmplPath string
}

// AddQr (de momento sin implementación activa; deja tu lógica cuando lo uses)
func (q *qrHandler) AddQr(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// GetAllQrs implements QrHandler.
func (q *qrHandler) GetAllQrs(w http.ResponseWriter, r *http.Request) {
	qrs, err := q.uc.GetAll()
	if err != nil {
		q.log.Printf("Error getting all QRs: %v", err)
		http.Error(w, "Error getting all QRs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// No llamamos WriteHeader(200); el primer Write fija 200
	_ = json.NewEncoder(w).Encode(qrs)
}

// GetQrById implements QrHandler.
func (q *qrHandler) GetQrById(id int, w http.ResponseWriter, r *http.Request) {
	qr, err := q.uc.GetById(id)
	if err != nil {
		q.log.Printf("Error getting QR by ID: %v", err)
		http.Error(w, "Error getting QR by ID", http.StatusInternalServerError)
		return
	}
	if qr == nil {
		http.Error(w, "QR not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(qr)
}

func (q *qrHandler) ClaimQr(id int, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	qr, err := q.uc.ClaimQr(id)
	if err != nil {
		if errors.Is(err, usecases.ErrQrFull) {
			w.WriteHeader(http.StatusForbidden) // o 409
			pretty, _ := json.MarshalIndent(qr, "", "  ")
			RenderTemplate(w, "denied.html", map[string]any{
				"message": "Se alcanzó el límite máximo de participantes para este QR.",
				"qr":      qr,
				"qrJSON":  string(pretty),
			})
			return
		}
		q.log.Printf("ClaimQr(%d) error: %v", id, err)
		http.Error(w, "Error claiming QR", http.StatusInternalServerError)
		return
	}

	pretty, _ := json.MarshalIndent(qr, "", "  ")
	RenderTemplate(w, "success.html", map[string]any{
		"qr":      qr,
		"qrJSON":  string(pretty),
		"message": "Asistencia registrada correctamente.",
	})
}
func NewQrHandler(logger *log.Logger, uc usecases.QrUseCase, tmplPath string) QrHandler {
	return &qrHandler{logger, uc, tmplPath}
}
