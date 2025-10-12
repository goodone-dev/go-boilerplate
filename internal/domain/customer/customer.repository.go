package customer

import (
	"github.com/goodonedev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICustomerRepository interface {
	database.IBaseRepository[gorm.DB, uuid.UUID, Customer]
}
