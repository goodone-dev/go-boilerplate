package mysql

import (
	"context"

	database "github.com/BagusAK95/go-skeleton/internal/infrastructure/database/sql"
	"github.com/BagusAK95/go-skeleton/internal/utils/tracer"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseRepo[M database.Entity] struct {
	Entity   M
	dbMaster *gorm.DB
	dbSlave  *gorm.DB
}

func NewBaseRepo[M database.Entity](dbConn dbConnection) database.IBaseRepository[M] {
	return &BaseRepo[M]{
		dbMaster: dbConn.Master,
		dbSlave:  dbConn.Slave,
	}
}

func (r *BaseRepo[M]) MasterConn() *gorm.DB {
	return r.dbMaster
}

func (r *BaseRepo[M]) SlaveConn() *gorm.DB {
	return r.dbSlave
}

func (r *BaseRepo[M]) GetAll(ctx context.Context) (res []M, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx)
	defer func() {
		span.EndSpan(err, res)
	}()

	builder := sq.
		Select("*").
		From(r.Entity.TableName()).
		Where(sq.Eq{
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

func (r *BaseRepo[M]) GetByID(ctx context.Context, ID uuid.UUID) (res M, err error) {
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

func (r *BaseRepo[M]) GetByIDLockTx(ctx context.Context, ID uuid.UUID, trx *gorm.DB) (res M, err error) {
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

	err = trx.WithContext(ctx).Raw(qry, args...).Scan(&res).Error
	if err != nil {
		return
	}

	return
}

func (r *BaseRepo[M]) GetByIDs(ctx context.Context, IDs []uuid.UUID) (res []M, err error) {
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

func (r *BaseRepo[M]) Pagination(ctx context.Context, filter map[string]any, page int, limit int) (res database.Pagination[M], err error) {
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

	var models []M
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

func (r *BaseRepo[M]) Create(ctx context.Context, req M) (res M, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, req)
	defer func() {
		span.EndSpan(err, res)
	}()

	err = r.dbMaster.WithContext(ctx).Create(&req).Error
	if err != nil {
		return req, err
	}

	return req, err
}

func (r *BaseRepo[M]) CreateWithTx(ctx context.Context, req M, trx *gorm.DB) (res M, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, req)
	defer func() {
		span.EndSpan(err, res)
	}()

	err = trx.WithContext(ctx).Create(&req).Error
	if err != nil {
		return req, err
	}

	return req, nil
}

func (r *BaseRepo[M]) CreateBulk(ctx context.Context, req []M) (res []M, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, req)
	defer func() {
		span.EndSpan(err, res)
	}()

	err = r.dbMaster.WithContext(ctx).CreateInBatches(req, database.InsertBatchSize).Error
	if err != nil {
		return req, err
	}

	return req, nil
}

func (r *BaseRepo[M]) CreateBulkWithTx(ctx context.Context, req []M, trx *gorm.DB) (res []M, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, req)
	defer func() {
		span.EndSpan(err, res)
	}()

	err = trx.WithContext(ctx).CreateInBatches(req, database.InsertBatchSize).Error
	if err != nil {
		return req, err
	}

	return req, nil
}

func (r *BaseRepo[M]) Update(ctx context.Context, ID uuid.UUID, req M) (res M, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, ID, req)
	defer func() {
		span.EndSpan(err, res)
	}()

	err = r.dbMaster.WithContext(ctx).Model(&req).Where("id=?", ID).Updates(req).Scan(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *BaseRepo[M]) UpdateWithTx(ctx context.Context, ID uuid.UUID, req M, trx *gorm.DB) (res M, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, ID, req)
	defer func() {
		span.EndSpan(err, res)
	}()

	err = trx.WithContext(ctx).Model(&req).Where("id=?", ID).Updates(req).Scan(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *BaseRepo[M]) UpdateWithMap(ctx context.Context, ID uuid.UUID, req map[string]any) (res M, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, ID, req)
	defer func() {
		span.EndSpan(err, res)
	}()

	err = r.dbMaster.WithContext(ctx).Model(&res).Where("id=?", ID).Updates(req).Scan(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *BaseRepo[M]) UpdateWithMapTx(ctx context.Context, ID uuid.UUID, req map[string]any, trx *gorm.DB) (res M, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, ID, req)
	defer func() {
		span.EndSpan(err, res)
	}()

	err = trx.WithContext(ctx).Model(&res).Where("id=?", ID).Updates(req).Scan(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *BaseRepo[M]) UpdateBulk(ctx context.Context, IDs []uuid.UUID, req map[string]any) (err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, IDs, req)
	defer func() {
		span.EndSpan(err)
	}()

	err = r.dbMaster.WithContext(ctx).Model(&r.Entity).Where("id IN ?", IDs).Updates(req).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepo[M]) UpdateBulkWithTx(ctx context.Context, IDs []uuid.UUID, req map[string]any, trx *gorm.DB) (err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, IDs, req)
	defer func() {
		span.EndSpan(err)
	}()

	err = trx.WithContext(ctx).Model(&r.Entity).Where("id IN ?", IDs).Updates(req).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepo[M]) Delete(ctx context.Context, ID uuid.UUID) (err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, ID)
	defer func() {
		span.EndSpan(err)
	}()

	err = r.dbMaster.WithContext(ctx).Delete(&r.Entity, ID).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepo[M]) DeleteWithTx(ctx context.Context, ID uuid.UUID, trx *gorm.DB) (err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, ID)
	defer func() {
		span.EndSpan(err)
	}()

	err = trx.WithContext(ctx).Delete(&r.Entity, ID).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepo[M]) DeleteBulk(ctx context.Context, IDs []uuid.UUID) (err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, IDs)
	defer func() {
		span.EndSpan(err)
	}()

	err = r.dbMaster.WithContext(ctx).Delete(&r.Entity, IDs).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepo[M]) DeleteBulkWithTx(ctx context.Context, IDs []uuid.UUID, trx *gorm.DB) (err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, IDs)
	defer func() {
		span.EndSpan(err)
	}()

	err = trx.WithContext(ctx).Delete(&r.Entity, IDs).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepo[M]) Begin(ctx context.Context) *gorm.DB {
	return r.dbMaster.WithContext(ctx).Begin()
}

func (r *BaseRepo[M]) Rollback(trx *gorm.DB) *gorm.DB {
	return trx.Rollback()
}

func (r *BaseRepo[M]) Commit(trx *gorm.DB) *gorm.DB {
	return trx.Commit()
}
