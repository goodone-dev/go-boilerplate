package repository

import (
	"github.com/goodone-dev/go-boilerplate/internal/domain/employee"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeRepository struct {
	database.IBaseRepository[gorm.DB, uuid.UUID, employee.Employee]
}

func NewEmployeeRepository(baseRepo database.IBaseRepository[gorm.DB, uuid.UUID, employee.Employee]) employee.IEmployeeRepository {
	return &EmployeeRepository{
		baseRepo,
	}
}
