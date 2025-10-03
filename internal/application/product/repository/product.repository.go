package repository

import (
	"github.com/BagusAK95/go-boilerplate/internal/domain/product"
	database "github.com/BagusAK95/go-boilerplate/internal/infrastructure/database/sql"
)

type ProductRepository struct {
	database.IBaseRepository[product.Product]
}

func NewProductRepo(baseRepo database.IBaseRepository[product.Product]) product.IProductRepository {
	return &ProductRepository{
		baseRepo,
	}
}
