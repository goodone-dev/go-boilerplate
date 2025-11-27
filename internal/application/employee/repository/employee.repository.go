package repository

import (
	"github.com/goodone-dev/go-boilerplate/internal/domain/employee"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type employeeRepository struct {
	database.BaseRepository[gorm.DB, uuid.UUID, employee.Employee]
}

func NewEmployeeRepository(baseRepo database.BaseRepository[gorm.DB, uuid.UUID, employee.Employee]) employee.EmployeeRepository {
	return &employeeRepository{
		baseRepo,
	}
}
