package repository

import (
	"github.com/goodone-dev/go-boilerplate/internal/domain/customer"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type customerRepository struct {
	database.BaseRepository[gorm.DB, uuid.UUID, customer.Customer]
}

func NewCustomerRepository(baseRepo database.BaseRepository[gorm.DB, uuid.UUID, customer.Customer]) customer.CustomerRepository {
	return &customerRepository{
		baseRepo,
	}
}
