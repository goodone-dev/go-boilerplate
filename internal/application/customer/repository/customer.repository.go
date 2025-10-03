package repository

import (
	"github.com/BagusAK95/go-skeleton/internal/domain/customer"
	database "github.com/BagusAK95/go-skeleton/internal/infrastructure/database/sql"
)

type CustomerRepository struct {
	database.IBaseRepository[customer.Customer]
}

func NewCustomerRepo(baseRepo database.IBaseRepository[customer.Customer]) customer.ICustomerRepository {
	return &CustomerRepository{
		baseRepo,
	}
}
