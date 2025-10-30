package repository

import (
	"github.com/goodone-dev/go-boilerplate/internal/domain/customer"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	database.IBaseRepository[gorm.DB, uuid.UUID, customer.Customer]
}

func NewCustomerRepository(baseRepo database.IBaseRepository[gorm.DB, uuid.UUID, customer.Customer]) customer.ICustomerRepository {
	return &CustomerRepository{
		baseRepo,
	}
}
