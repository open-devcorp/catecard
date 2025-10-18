package handlers

import (
	"catecard/pkg/domain/entities"
	"catecard/pkg/domain/usecases"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type QrHandler interface {
	// Define QR handler methods here
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

// Añade los imports "math/rand" y "time"
func (q *qrHandler) AddQr(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	forum := rand.Intn(5) + 1

	qr := entities.Qr{
		Forum:   forum,
		GroupId: 1,
	}

	// get the authenticated user from request context (may be nil if not logged)
	user := GetUserFromRequest(r)
	err := q.uc.Add(user, &qr)
	if err != nil {
		q.log.Printf("Error adding QR: %v", err)
		http.Error(w, "Error adding QR", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(qrs)
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(qr)
}

// UpdateQr implements QrHandler.
func (q *qrHandler) ClaimQr(id int, w http.ResponseWriter, r *http.Request) {
	qr, err := q.uc.ClaimQr(id)
	if err != nil {
		q.log.Printf("Error claiming QR: %v", err)
		// render denied view with error message
		RenderTemplate(w, "denied.html", map[string]interface{}{
			"Message": err.Error(),
			"Data":    qr,
		})
		return
	}

	RenderTemplate(w, "success.html", map[string]interface{}{
		"message": "QR claimed successfully",
		"qr":      qr,
	})
}

func NewQrHandler(logger *log.Logger, uc usecases.QrUseCase, tmplPath string) QrHandler {
	return &qrHandler{logger, uc, tmplPath}
}
