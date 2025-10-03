package sql

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const InsertBatchSize = 100

type Pagination[M Entity] struct {
	Data    []M  `json:"data"`
	HasNext bool `json:"has_next"`
	HasPrev bool `json:"has_prev"`
}

type IBaseRepository[M Entity] interface {
	MasterConn() *gorm.DB
	SlaveConn() *gorm.DB
	GetAll(ctx context.Context) ([]M, error)
	GetByID(ctx context.Context, ID uuid.UUID) (M, error)
	GetByIDLockTx(ctx context.Context, ID uuid.UUID, trx *gorm.DB) (M, error)
	GetByIDs(ctx context.Context, IDs []uuid.UUID) ([]M, error)
	Pagination(ctx context.Context, filter map[string]any, page int, limit int) (res Pagination[M], err error)
	Create(ctx context.Context, model M) (M, error)
	CreateWithTx(ctx context.Context, model M, trx *gorm.DB) (M, error)
	CreateBulk(ctx context.Context, models []M) ([]M, error)
	CreateBulkWithTx(ctx context.Context, models []M, trx *gorm.DB) ([]M, error)
	Update(ctx context.Context, ID uuid.UUID, model M) (M, error)
	UpdateWithTx(ctx context.Context, ID uuid.UUID, model M, trx *gorm.DB) (M, error)
	UpdateWithMap(ctx context.Context, ID uuid.UUID, payload map[string]any) (M, error)
	UpdateWithMapTx(ctx context.Context, ID uuid.UUID, payload map[string]any, trx *gorm.DB) (M, error)
	UpdateBulk(ctx context.Context, IDs []uuid.UUID, payload map[string]any) error
	UpdateBulkWithTx(ctx context.Context, IDs []uuid.UUID, payload map[string]any, trx *gorm.DB) error
	Delete(ctx context.Context, ID uuid.UUID) error
	DeleteWithTx(ctx context.Context, ID uuid.UUID, trx *gorm.DB) error
	DeleteBulk(ctx context.Context, IDs []uuid.UUID) error
	DeleteBulkWithTx(ctx context.Context, IDs []uuid.UUID, trx *gorm.DB) error
	Begin(ctx context.Context) *gorm.DB
	Rollback(trx *gorm.DB) *gorm.DB
	Commit(trx *gorm.DB) *gorm.DB
}
