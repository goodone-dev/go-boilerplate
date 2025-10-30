package postgres

import (
	"context"
	"math"

	sq "github.com/Masterminds/squirrel"
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
	"gorm.io/gorm"
)

type BaseRepo[D any, I any, E database.Entity] struct {
	Entity   E
	dbMaster *gorm.DB
	dbSlave  *gorm.DB
}

func NewBaseRepository[D any, I any, E database.Entity](dbConn postgresConnection) database.IBaseRepository[D, I, E] {
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

func (r *BaseRepo[D, I, E]) FindAll(ctx context.Context, filter map[string]any) (res []E, err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, filter)
	defer func() {
		span.Stop(err, res)
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
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, ID)
	defer func() {
		span.Stop(err, res)
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
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, ID)
	defer func() {
		span.Stop(err, res)
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
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, IDs)
	defer func() {
		span.Stop(err, res)
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

func (r *BaseRepo[D, I, E]) FindByOffset(ctx context.Context, filter map[string]any, sort []string, size int, page int) (res database.Pagination[E], err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, filter, sort, size, page)
	defer func() {
		span.Stop(err, res)
	}()

	if size <= 0 {
		size = 10
	}
	if page <= 0 {
		page = 1
	}

	filter["deleted_at"] = nil

	builder := sq.
		Select("COUNT(*)").
		From(r.Entity.TableName()).
		Where(filter)

	qry, args, err := builder.ToSql()
	if err != nil {
		return
	}

	var total int64
	err = r.dbSlave.WithContext(ctx).Raw(qry, args...).Scan(&total).Error
	if err != nil {
		return
	}

	builder = sq.
		Select("*").
		From(r.Entity.TableName()).
		Where(filter).
		OrderBy(sort...).
		Limit(uint64(size)).
		Offset(uint64((page - 1) * size))

	qry, args, err = builder.ToSql()
	if err != nil {
		return
	}

	var models []E
	err = r.dbSlave.WithContext(ctx).Raw(qry, args...).Scan(&models).Error
	if err != nil {
		return
	}

	var pages int
	if total > 0 {
		pages = int(math.Ceil(float64(total) / float64(size)))
	}

	res.Data = models
	res.Metadata.Total = &total
	res.Metadata.Pages = &pages
	res.Metadata.Page = &page
	res.Metadata.Size = &size

	return
}

func (r *BaseRepo[D, I, E]) FindByCursor(ctx context.Context, filter map[string]any, sort []string, size int, next *I) (res database.Pagination[E], err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, filter, sort, size, next)
	defer func() {
		span.Stop(err, res)
	}()

	if size <= 0 {
		size = 10
	}

	filter["deleted_at"] = nil
	if next != nil {
		filter["id > ?"] = *next
	}

	builder := sq.
		Select("COUNT(*)").
		From(r.Entity.TableName()).
		Where(filter)

	qry, args, err := builder.ToSql()
	if err != nil {
		return
	}

	var total int64
	err = r.dbSlave.WithContext(ctx).Raw(qry, args...).Scan(&total).Error
	if err != nil {
		return
	}

	builder = sq.
		Select("*").
		From(r.Entity.TableName()).
		Where(filter).
		OrderBy(sort...).
		Limit(uint64(size))

	qry, args, err = builder.ToSql()
	if err != nil {
		return
	}

	var models []E
	err = r.dbSlave.WithContext(ctx).Raw(qry, args...).Scan(&models).Error
	if err != nil {
		return
	}

	var pages int
	if total > 0 {
		pages = int(math.Ceil(float64(total) / float64(size)))
	}

	res.Data = models
	res.Metadata.Total = &total
	res.Metadata.Pages = &pages
	res.Metadata.Size = &size

	return
}

// TODO: Check 'res' is still necessary
func (r *BaseRepo[D, I, E]) Insert(ctx context.Context, req E, trx *D) (res E, err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, req)
	defer func() {
		span.Stop(err, res)
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

// TODO: Check 'res' is still necessary
func (r *BaseRepo[D, I, E]) InsertMany(ctx context.Context, req []E, trx *D) (res []E, err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, req)
	defer func() {
		span.Stop(err, res)
	}()

	db := r.dbMaster
	if trx, ok := any(trx).(*gorm.DB); ok {
		db = trx
	}

	err = db.WithContext(ctx).CreateInBatches(req, config.InsertBatchSize).Error
	if err != nil {
		return req, err
	}

	return req, nil
}

func (r *BaseRepo[D, I, E]) Update(ctx context.Context, req E, trx *D) (err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, req)
	defer func() {
		span.Stop(err)
	}()

	db := r.dbMaster
	if trx, ok := any(trx).(*gorm.DB); ok {
		db = trx
	}

	err = db.WithContext(ctx).Save(&req).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepo[D, I, E]) UpdateById(ctx context.Context, ID I, req map[string]any, trx *D) (res E, err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, ID, req)
	defer func() {
		span.Stop(err, res)
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

// TODO: Check 'res' is needed
func (r *BaseRepo[D, I, E]) UpdateByIds(ctx context.Context, IDs []I, req map[string]any, trx *D) (err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, IDs, req)
	defer func() {
		span.Stop(err)
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

// TODO: Check 'res' is needed
func (r *BaseRepo[D, I, E]) UpdateMany(ctx context.Context, filter map[string]any, req map[string]any, trx *D) (err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, filter, req)
	defer func() {
		span.Stop(err)
	}()

	db := r.dbMaster
	if trx, ok := any(trx).(*gorm.DB); ok {
		db = trx
	}

	err = db.WithContext(ctx).Model(&r.Entity).Where(filter).Updates(req).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepo[D, I, E]) DeleteById(ctx context.Context, ID I, trx *D) (err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, ID)
	defer func() {
		span.Stop(err)
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
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, IDs)
	defer func() {
		span.Stop(err)
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

func (r *BaseRepo[D, I, E]) DeleteMany(ctx context.Context, filter map[string]any, trx *D) (err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, filter)
	defer func() {
		span.Stop(err)
	}()

	db := r.dbMaster
	if trx, ok := any(trx).(*gorm.DB); ok {
		db = trx
	}

	err = db.WithContext(ctx).Delete(&r.Entity, filter).Error
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
