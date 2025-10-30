package repository

import (
	"github.com/goodone-dev/go-boilerplate/internal/domain/product"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	database.IBaseRepository[gorm.DB, uuid.UUID, product.Product]
}

func NewProductRepository(baseRepo database.IBaseRepository[gorm.DB, uuid.UUID, product.Product]) product.IProductRepository {
	return &productRepository{
		baseRepo,
	}
}
