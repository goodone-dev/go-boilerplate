package repository

import (
	"github.com/BagusAK95/go-skeleton/internal/domain/product"
	database "github.com/BagusAK95/go-skeleton/internal/infrastructure/database/sql"
)

type productRepo struct {
	database.IBaseRepository[product.Product]
}

func NewProductRepo(baseRepo database.IBaseRepository[product.Product]) product.IProductRepository {
	return &productRepo{
		baseRepo,
	}
}
