package employee

import (
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
)

type Employee struct {
	database.BaseEntity[uuid.UUID] `bson:",inline"`
	Name                           string `json:"name" bson:"name"`
	Email                          string `json:"email" bson:"email"`
	Role                           string `json:"role" bson:"role"`
}

func (Employee) TableName() string {
	return "employees"
}

func (Employee) RepositoryName() string {
	return "EmployeeRepository"
}
