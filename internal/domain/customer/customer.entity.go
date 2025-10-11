package customer

import (
	"github.com/BagusAK95/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
)

type Customer struct {
	database.BaseEntity[uuid.UUID]
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (Customer) TableName() string {
	return "customers"
}

func (Customer) RepositoryName() string {
	return "CustomerRepository"
}
