package employee

import (
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IEmployeeRepository interface {
	database.IBaseRepository[gorm.DB, uuid.UUID, Employee]
}
