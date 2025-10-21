package handlers

import (
	"catecard/pkg/domain/entities"
	"catecard/pkg/domain/usecases"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type CatechumenHandler interface {
	AddCatechumen(w http.ResponseWriter, r *http.Request)
	UpdateCatechumen(w http.ResponseWriter, r *http.Request)
	GetAllCatechumens(w http.ResponseWriter, r *http.Request)
	GetCatechumenById(id int, w http.ResponseWriter, r *http.Request)
}

type catechumenHandler struct {
	log      *log.Logger
	uc       usecases.CatechumenUseCase
	tmplPath string
}

func NewCatechumenHandler(logger *log.Logger, uc usecases.CatechumenUseCase, tmplPath string) CatechumenHandler {
	return &catechumenHandler{logger, uc, tmplPath}
}

// AddCatechumen implements CatechumenHandler.
func (c *catechumenHandler) AddCatechumen(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.log.Printf("Error parsing form: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	fullName := r.FormValue("name")
	age := r.FormValue("age")
	groupIdStr := r.FormValue("group_id")

	c.log.Printf("group_id: %s", groupIdStr)

	if groupIdStr == "" {
		c.log.Printf("group_id is empty")
		writeJSONError(w, http.StatusBadRequest, "group_id is required")
		return
	}

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		c.log.Printf("Error parsing group ID: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid group ID value")
		return
	}

	catechumen := entities.NewCatechumen(fullName, age, groupId)

	addedCatechumen, qr, err := c.uc.Add(&catechumen)
	if err != nil {
		c.log.Printf("Error adding catechumen: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error adding catechumen")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := map[string]interface{}{
		"status":     "ok",
		"message":    "Catechumen created",
		"catechumen": addedCatechumen,
		"qr":         qr,
	}

	json.NewEncoder(w).Encode(resp)
}

// GetAllCatechumens implements CatechumenHandler.
func (c *catechumenHandler) GetAllCatechumens(w http.ResponseWriter, r *http.Request) {

	catechumens, err := c.uc.GetAll()
	if err != nil {
		c.log.Printf("Error retrieving catechumens: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error retrieving catechumens")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]interface{}{
		"status":      "ok",
		"catechumens": catechumens,
	}
	json.NewEncoder(w).Encode(resp)
}

// GetCatechumenById implements CatechumenHandler.
func (c *catechumenHandler) GetCatechumenById(id int, w http.ResponseWriter, r *http.Request) {
	catechumen, err := c.uc.GetById(id)
	if err != nil {
		c.log.Printf("Error retrieving catechumen by ID: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error retrieving catechumen")
		return
	}
	if catechumen == nil {
		writeJSONError(w, http.StatusNotFound, "Catechumen not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]interface{}{
		"status":     "ok",
		"catechumen": catechumen,
	}
	json.NewEncoder(w).Encode(resp)

}

// UpdateCatechumen implements CatechumenHandler.
func (c *catechumenHandler) UpdateCatechumen(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.log.Printf("Error parsing form: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	idStr := r.FormValue("id")
	fullName := r.FormValue("name")
	age := r.FormValue("age")
	groupIdStr := r.FormValue("group_id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.log.Printf("Error parsing catechumen ID: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid catechumen ID value")
		return
	}

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		c.log.Printf("Error parsing group ID: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid group ID value")
		return
	}

	catechumen := &entities.Catechumen{
		ID:       id,
		FullName: fullName,
		Age:      age,
		GroupId:  groupId,
	}

	updatedCatechumen, err := c.uc.Update(catechumen)
	if err != nil {
		c.log.Printf("Error updating catechumen: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Error updating catechumen")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]interface{}{
		"status":     "ok",
		"message":    "Catechumen updated",
		"catechumen": updatedCatechumen,
	}

	json.NewEncoder(w).Encode(resp)
}
