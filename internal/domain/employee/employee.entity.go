package employee

import database "github.com/BagusAK95/go-boilerplate/internal/infrastructure/database/sql"

type Employee struct {
	database.BaseEntity
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
