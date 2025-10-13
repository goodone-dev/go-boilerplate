package database

import (
	"context"
)

const InsertBatchSize = 100

type Pagination[E Entity] struct {
	Data    []E  `json:"data"`
	HasNext bool `json:"has_next"`
	HasPrev bool `json:"has_prev"`
}

type IBaseRepository[D any, I any, E Entity] interface {
	// Database
	MasterDB() *D
	SlaveDB() *D

	// Common Query
	Find(ctx context.Context, filter map[string]any) ([]E, error)
	FindWithPagination(ctx context.Context, filter map[string]any, page int, size int) (res Pagination[E], err error)
	FindById(ctx context.Context, ID I) (*E, error)
	FindByIdAndLock(ctx context.Context, ID I, trx *D) (*E, error)
	FindByIds(ctx context.Context, IDs []I) ([]E, error)
	Insert(ctx context.Context, model E, trx *D) (E, error)
	InsertMany(ctx context.Context, models []E, trx *D) ([]E, error)
	UpdateById(ctx context.Context, ID I, payload map[string]any, trx *D) (E, error)
	UpdateByIds(ctx context.Context, IDs []I, payload map[string]any, trx *D) error
	DeleteById(ctx context.Context, ID I, trx *D) error
	DeleteByIds(ctx context.Context, IDs []I, trx *D) error

	// Transaction
	Begin(ctx context.Context) (*D, error)
	Rollback(trx *D) *D
	Commit(trx *D) *D
}
