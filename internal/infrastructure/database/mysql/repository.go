package mysql

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/goodonedev/go-boilerplate/internal/infrastructure/database"
	"github.com/goodonedev/go-boilerplate/internal/utils/tracer"
	"gorm.io/gorm"
)

type BaseRepo[D any, I any, E database.Entity] struct {
	Entity   E
	dbMaster *gorm.DB
	dbSlave  *gorm.DB
}

func NewBaseRepo[D any, I any, E database.Entity](dbConn mysqlConnection) database.IBaseRepository[D, I, E] {
	return &BaseRepo[D, I, E]{
		dbMaster: dbConn.Master,
		dbSlave:  dbConn.Slave,
	}
}

func (r *BaseRepo[D, I, E]) MasterDB() *D {
	return any(r.dbMaster).(*D)
}

func (r *BaseRepo[D, I, E]) SlaveDB() *D {
	return any(r.dbSlave).(*D)
}

func (r *BaseRepo[D, I, E]) Find(ctx context.Context, filter map[string]any) (res []E, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, filter)
	defer func() {
		span.EndSpan(err, res)
	}()

	filter["deleted_at"] = nil

	builder := sq.
		Select("*").
		From(r.Entity.TableName()).
		Where(filter)

	qry, args, err := builder.ToSql()
	if err != nil {
		return
	}

	err = r.dbSlave.WithContext(ctx).Raw(qry, args...).Scan(&res).Error
	if err != nil {
		return
	}

	return
}

func (r *BaseRepo[D, I, E]) FindById(ctx context.Context, ID I) (res *E, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, ID)
	defer func() {
		span.EndSpan(err, res)
	}()

	builder := sq.
		Select("*").
		From(r.Entity.TableName()).
		Where(sq.Eq{
			"id":         ID,
			"deleted_at": nil,
		})

	qry, args, err := builder.ToSql()
	if err != nil {
		return
	}

	err = r.dbSlave.WithContext(ctx).Raw(qry, args...).Scan(&res).Error
	if err != nil {
		return
	}

	return
}

func (r *BaseRepo[D, I, E]) FindByIdAndLock(ctx context.Context, ID I, trx *D) (res *E, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, ID)
	defer func() {
		span.EndSpan(err, res)
	}()

	builder := sq.
		Select("*").
		From(r.Entity.TableName()).
		Where(sq.Eq{
			"id":         ID,
			"deleted_at": nil,
		}).
		Suffix("FOR UPDATE")

	qry, args, err := builder.ToSql()
	if err != nil {
		return
	}

	db := r.dbSlave
	if trx, ok := any(trx).(*gorm.DB); ok {
		db = trx
	}

	err = db.WithContext(ctx).Raw(qry, args...).Scan(&res).Error
	if err != nil {
		return
	}

	return
}

func (r *BaseRepo[D, I, E]) FindByIds(ctx context.Context, IDs []I) (res []E, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, IDs)
	defer func() {
		span.EndSpan(err, res)
	}()

	builder := sq.
		Select("*").
		From(r.Entity.TableName()).
		Where(sq.Eq{
			"id":         IDs,
			"deleted_at": nil,
		})

	qry, args, err := builder.ToSql()
	if err != nil {
		return
	}

	err = r.dbSlave.WithContext(ctx).Raw(qry, args...).Scan(&res).Error
	if err != nil {
		return
	}

	return
}

func (r *BaseRepo[D, I, E]) FindWithPagination(ctx context.Context, filter map[string]any, page int, limit int) (res database.Pagination[E], err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, filter, page, limit)
	defer func() {
		span.EndSpan(err, res)
	}()

	filter["deleted_at"] = nil

	builder := sq.
		Select("*").
		From(r.Entity.TableName()).
		Where(filter).
		OrderBy("id DESC").
		Limit(uint64(limit + 1)).
		Offset(uint64((page - 1) * limit))

	qry, args, err := builder.ToSql()
	if err != nil {
		return
	}

	var models []E
	err = r.dbSlave.WithContext(ctx).Raw(qry, args...).Scan(&models).Error
	if err != nil {
		return
	}

	if len(models) > limit {
		res.HasNext = true
		models = models[:limit]
	}

	if page > 1 {
		res.HasPrev = true
	}

	res.Data = models

	return
}

func (r *BaseRepo[D, I, E]) Insert(ctx context.Context, req E, trx *D) (res E, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, req)
	defer func() {
		span.EndSpan(err, res)
	}()

	db := r.dbMaster
	if trx, ok := any(trx).(*gorm.DB); ok {
		db = trx
	}

	err = db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return req, err
	}

	return req, nil
}

func (r *BaseRepo[D, I, E]) InsertMany(ctx context.Context, req []E, trx *D) (res []E, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, req)
	defer func() {
		span.EndSpan(err, res)
	}()

	db := r.dbMaster
	if trx, ok := any(trx).(*gorm.DB); ok {
		db = trx
	}

	err = db.WithContext(ctx).CreateInBatches(req, database.InsertBatchSize).Error
	if err != nil {
		return req, err
	}

	return req, nil
}

func (r *BaseRepo[D, I, E]) UpdateById(ctx context.Context, ID I, req map[string]any, trx *D) (res E, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, ID, req)
	defer func() {
		span.EndSpan(err, res)
	}()

	db := r.dbMaster
	if trx, ok := any(trx).(*gorm.DB); ok {
		db = trx
	}

	err = db.WithContext(ctx).Model(&res).Where("id=?", ID).Updates(req).Scan(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *BaseRepo[D, I, E]) UpdateByIds(ctx context.Context, IDs []I, req map[string]any, trx *D) (err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, IDs, req)
	defer func() {
		span.EndSpan(err)
	}()

	db := r.dbMaster
	if trx, ok := any(trx).(*gorm.DB); ok {
		db = trx
	}

	err = db.WithContext(ctx).Model(&r.Entity).Where("id IN ?", IDs).Updates(req).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepo[D, I, E]) DeleteById(ctx context.Context, ID I, trx *D) (err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, ID)
	defer func() {
		span.EndSpan(err)
	}()

	db := r.dbMaster
	if trx, ok := any(trx).(*gorm.DB); ok {
		db = trx
	}

	err = db.WithContext(ctx).Delete(&r.Entity, ID).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepo[D, I, E]) DeleteByIds(ctx context.Context, IDs []I, trx *D) (err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, IDs)
	defer func() {
		span.EndSpan(err)
	}()

	db := r.dbMaster
	if trx, ok := any(trx).(*gorm.DB); ok {
		db = trx
	}

	err = db.WithContext(ctx).Delete(&r.Entity, IDs).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepo[D, I, E]) Begin(ctx context.Context) (trx *D, err error) {
	db := r.dbMaster.WithContext(ctx).Begin()
	if db.Error != nil {
		return trx, db.Error
	}

	return any(db).(*D), nil
}

func (r *BaseRepo[D, I, E]) Rollback(trx *D) *D {
	db, ok := any(trx).(*gorm.DB)
	if !ok {
		return trx
	}

	db = db.Rollback()

	return any(db).(*D)
}

func (r *BaseRepo[D, I, E]) Commit(trx *D) *D {
	db, ok := any(trx).(*gorm.DB)
	if !ok {
		return trx
	}

	db = db.Commit()

	return any(db).(*D)
}
