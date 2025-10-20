package database

import (
	"context"
)

type Pagination[E Entity] struct {
	Data       []E                   `json:"data"`
	Metadata   *PaginationMetadata   `json:"metadata,omitempty"`
	Navigation *PaginationNavigation `json:"navigation,omitempty"`
}

type PaginationMetadata struct {
	Total *int64 `json:"total,omitempty"`
	Pages *int   `json:"pages,omitempty"`
	Page  *int   `json:"page,omitempty"`
	Size  *int   `json:"size,omitempty"`
}

type PaginationNavigation struct {
	First *string `json:"first,omitempty"`
	Next  *string `json:"next,omitempty"`
	Prev  *string `json:"prev,omitempty"`
	Last  *string `json:"last,omitempty"`
}

type IBaseRepository[D any, I any, E Entity] interface {
	MasterDB() *D
	SlaveDB() *D

	FindAll(ctx context.Context, filter map[string]any) ([]E, error)
	FindById(ctx context.Context, ID I) (*E, error)
	FindByIdAndLock(ctx context.Context, ID I, trx *D) (*E, error)
	FindByIds(ctx context.Context, IDs []I) ([]E, error)
	FindByOffset(ctx context.Context, filter map[string]any, sort []string, size int, page int) (res Pagination[E], err error)
	FindByCursor(ctx context.Context, filter map[string]any, sort []string, size int, next *I) (res Pagination[E], err error)

	Insert(ctx context.Context, model E, trx *D) (E, error)
	InsertMany(ctx context.Context, models []E, trx *D) ([]E, error)

	Update(ctx context.Context, model E, trx *D) error
	UpdateById(ctx context.Context, ID I, payload map[string]any, trx *D) (E, error)
	UpdateByIds(ctx context.Context, IDs []I, payload map[string]any, trx *D) error
	UpdateMany(ctx context.Context, filter map[string]any, payload map[string]any, trx *D) error

	DeleteById(ctx context.Context, ID I, trx *D) error
	DeleteByIds(ctx context.Context, IDs []I, trx *D) error
	DeleteMany(ctx context.Context, filter map[string]any, trx *D) error

	Begin(ctx context.Context) (*D, error)
	Rollback(trx *D) *D
	Commit(trx *D) *D
}
