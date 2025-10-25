package entities

import "time"

type ScanCatechumen struct {
	ID           int    `json:"id"`
	CatechumenID int    `json:"catechumen_id"`
	ScanID       int    `json:"scan_id"`
	CreatedAt    string `json:"created_at"`
}

func NewScanCatechumen(catechumenID int, scanID int) *ScanCatechumen {
	return &ScanCatechumen{
		CatechumenID: catechumenID,
		ScanID:       scanID,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}
}
