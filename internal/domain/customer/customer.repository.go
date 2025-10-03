package customer

import database "github.com/BagusAK95/go-boilerplate/internal/infrastructure/database/sql"

type ICustomerRepository interface {
	database.IBaseRepository[Customer]
}
