package product

import database "github.com/BagusAK95/go-skeleton/internal/infrastructure/database/sql"

type Product struct {
	database.BaseEntity
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func (Product) TableName() string {
	return "products"
}

func (Product) RepositoryName() string {
	return "ProductRepository"
}
