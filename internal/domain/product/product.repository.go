package product

import database "github.com/BagusAK95/go-boilerplate/internal/infrastructure/database/sql"

type IProductRepository interface {
	database.IBaseRepository[Product]
}
