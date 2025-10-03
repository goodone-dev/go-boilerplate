package order

import database "github.com/BagusAK95/go-skeleton/internal/infrastructure/database/sql"

type IOrderItemRepository interface {
	database.IBaseRepository[OrderItem]
}
