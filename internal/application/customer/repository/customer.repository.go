package repository

import (
	"github.com/BagusAK95/go-boilerplate/internal/domain/customer"
	database "github.com/BagusAK95/go-boilerplate/internal/infrastructure/database/sql"
)

type CustomerRepository struct {
	database.IBaseRepository[customer.Customer]
}

func NewCustomerRepo(baseRepo database.IBaseRepository[customer.Customer]) customer.ICustomerRepository {
	return &CustomerRepository{
		baseRepo,
	}
}
