package customer

import (
	database "github.com/BagusAK95/go-boilerplate/internal/infrastructure/database/sql"
)

type Customer struct {
	database.BaseEntity
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (Customer) TableName() string {
	return "customers"
}

func (Customer) RepositoryName() string {
	return "CustomerRepository"
}
