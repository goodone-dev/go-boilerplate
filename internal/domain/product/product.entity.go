package product

import (
	"github.com/BagusAK95/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
)

type Product struct {
	database.BaseEntity[uuid.UUID] `bson:",inline"`
	Name                           string  `json:"name" bson:"name"`
	Description                    string  `json:"description" bson:"description"`
	Price                          float64 `json:"price" bson:"price"`
	Stock                          int     `json:"stock" bson:"stock"`
}

func (Product) TableName() string {
	return "products"
}

func (Product) RepositoryName() string {
	return "ProductRepository"
}
