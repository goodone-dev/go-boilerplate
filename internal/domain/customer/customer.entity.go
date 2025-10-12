package customer

import (
	"github.com/goodonedev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
)

type Customer struct {
	database.BaseEntity[uuid.UUID] `bson:",inline"`
	Name                           string `json:"name" bson:"name"`
	Email                          string `json:"email" bson:"email"`
}

func (Customer) TableName() string {
	return "customers"
}

func (Customer) RepositoryName() string {
	return "CustomerRepository"
}
