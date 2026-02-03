package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/goodone-dev/go-boilerplate/internal/domain/order"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderItemRepository struct {
	database.BaseRepository[gorm.DB, uuid.UUID, order.OrderItem]
}

func NewOrderItemRepository(baseRepo database.BaseRepository[gorm.DB, uuid.UUID, order.OrderItem]) order.OrderItemRepository {
	return &orderItemRepository{
		baseRepo,
	}
}

func (r *orderItemRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) (res []order.DetailOrderItem, err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"orderID": orderID,
	})

	defer func() {
		span.SetFunctionOutput(tracer.Metadata{
			"response": res,
		}).End(err)
	}()

	builder := sq.
		Select("order_items.*, products.name as product_name").
		From("order_items").
		LeftJoin("products ON products.id = order_items.product_id").
		Where(sq.Eq{
			"order_items.order_id":   orderID,
			"order_items.deleted_at": nil,
		})

	qry, args, err := builder.ToSql()
	if err != nil {
		return
	}

	err = r.SlaveDB().WithContext(ctx).Raw(qry, args...).Scan(&res).Error
	if err != nil {
		return
	}

	return
}
