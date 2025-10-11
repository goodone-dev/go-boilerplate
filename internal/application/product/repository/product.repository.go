package repository

import (
	"github.com/BagusAK95/go-boilerplate/internal/domain/product"
	"github.com/BagusAK95/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository struct {
	database.IBaseRepository[gorm.DB, uuid.UUID, product.Product]
}

func NewProductRepo(baseRepo database.IBaseRepository[gorm.DB, uuid.UUID, product.Product]) product.IProductRepository {
	return &ProductRepository{
		baseRepo,
	}
}
