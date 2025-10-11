package product

import (
	"github.com/BagusAK95/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
)

type Product struct {
	database.BaseEntity[uuid.UUID]
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
