package order

import database "github.com/BagusAK95/go-boilerplate/internal/infrastructure/database/sql"

type IOrderItemRepository interface {
	database.IBaseRepository[OrderItem]
}
