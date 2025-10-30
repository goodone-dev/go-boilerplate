package repository

import (
	"github.com/goodone-dev/go-boilerplate/internal/domain/product"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	database.BaseRepository[gorm.DB, uuid.UUID, product.Product]
}

func NewProductRepository(baseRepo database.BaseRepository[gorm.DB, uuid.UUID, product.Product]) product.ProductRepository {
	return &productRepository{
		baseRepo,
	}
}
