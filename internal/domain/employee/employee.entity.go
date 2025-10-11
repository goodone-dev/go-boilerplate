package employee

import (
	"github.com/BagusAK95/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
)

type Employee struct {
	database.BaseEntity[uuid.UUID]
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (Employee) TableName() string {
	return "employees"
}

func (Employee) RepositoryName() string {
	return "EmployeeRepository"
}
