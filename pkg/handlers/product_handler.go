package handlers

import (
	"catecard/pkg/domain/entities"
	"catecard/pkg/domain/usecases"
	"log"
	"net/http"
	"strconv"
)

type ProductHandler interface {
	AddProduct(w http.ResponseWriter, r *http.Request)
	GetAllProducts(w http.ResponseWriter, r *http.Request)
}

type productHandler struct {
	log      *log.Logger
	uc       usecases.ProductUseCase
	tmplPath string
}

// AddProduct implements ProductHandler.
func (p *productHandler) AddProduct(w http.ResponseWriter, r *http.Request) {

	// form
	if err := r.ParseForm(); err != nil {
		p.log.Printf("Error parsing form: %v", err)
		http.Error(w, "invalid form data", http.StatusBadRequest)

		return

	}

	name := r.FormValue("name")

	price := r.FormValue("price")

	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		p.log.Printf("Error parsing price: %v", err)
		http.Error(w, "invalid price value", http.StatusBadRequest)
		return
	}

	product := entities.NewProduct(name, float32(priceFloat))

	err = p.uc.Add(&product)

	if err == nil {
		p.log.Printf("Product added")
	}

	http.Redirect(w, r, "/products", http.StatusSeeOther)
}

// GetAllProducts implements ProductHandler.
func (p *productHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {

	products, err := p.uc.GetAll()

	if err != nil {

		p.log.Printf("Error getting products: %v", err)
		http.Error(w, "Error retrieving products", http.StatusInternalServerError)
		return
	}
	RenderTemplate(w, "products.html", products)

}

func NewProductHandler(logger *log.Logger, uc usecases.ProductUseCase, tmplPath string) ProductHandler {

	return &productHandler{logger, uc, tmplPath}

}
