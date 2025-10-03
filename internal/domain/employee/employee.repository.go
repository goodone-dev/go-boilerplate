package employee

import database "github.com/BagusAK95/go-boilerplate/internal/infrastructure/database/sql"

type IEmployeeRepository interface {
	database.IBaseRepository[Employee]
}
