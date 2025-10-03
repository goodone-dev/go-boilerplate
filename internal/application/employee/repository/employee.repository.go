package repository

import (
	"github.com/BagusAK95/go-skeleton/internal/domain/employee"
	database "github.com/BagusAK95/go-skeleton/internal/infrastructure/database/sql"
)

type employeeRepo struct {
	database.IBaseRepository[employee.Employee]
}

func NewEmployeeRepo(baseRepo database.IBaseRepository[employee.Employee]) employee.IEmployeeRepository {
	return &employeeRepo{
		baseRepo,
	}
}
