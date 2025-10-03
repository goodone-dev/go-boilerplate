package employee

import database "github.com/BagusAK95/go-skeleton/internal/infrastructure/database/sql"

type IEmployeeRepository interface {
	database.IBaseRepository[Employee]
}
