package repositories

import (
	"catecard/pkg/domain/entities"
	"database/sql"
	"log"
)

type productRepository struct {
	log *log.Logger
	db  *sql.DB
}

// Add implements ProductRepository.
func (p *productRepository) Add(product *entities.Product) error {

	query := `INSERT INTO products(name, price) VALUES (?,?)`
	result, err := p.db.Exec(query, product.Name, product.Price)
	if err != nil {
		p.log.Printf("Error inserting product: %v", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		p.log.Printf("Error getting last inser ID: %v", err)
		return err
	}
	product.ID = int(id)
	return nil
}

// GetAll implements ProductRepository.
func (p *productRepository) GetAll() ([]*entities.Product, error) {

	query := `SELECT id, name, price FROM products`
	rows, err := p.db.Query(query)

	if err != nil {
		p.log.Printf("Error querying products: %v", err)
		return nil, err
	}
	defer rows.Close()
	var products []*entities.Product
	for rows.Next() {
		product := &entities.Product{}
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			p.log.Printf("Error scanning product: %v", err)
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil

}

func NewProductRepository(logger *log.Logger, db *sql.DB) ProductRepository {
	return &productRepository{logger, db}
}
