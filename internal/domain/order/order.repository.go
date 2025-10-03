package order

import (
	database "github.com/BagusAK95/go-boilerplate/internal/infrastructure/database/sql"
)

type IOrderRepository interface {
	database.IBaseRepository[Order]
}
