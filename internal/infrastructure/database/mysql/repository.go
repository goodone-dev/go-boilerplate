package mysql

import (
	"context"
	"math"

	sq "github.com/Masterminds/squirrel"
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
	"gorm.io/gorm"
)

type baseRepo[D any, I any, E database.Entity] struct {
	Entity   E
	dbMaster *gorm.DB
	dbSlave  *gorm.DB
}

func NewBaseRepository[D any, I any, E database.Entity](dbConn *mysqlConnection) database.BaseRepository[D, I, E] {
	return &baseRepo[D, I, E]{
		dbMaster: dbConn.Master,
		dbSlave:  dbConn.Slave,
	}
}

func (r *baseRepo[D, I, E]) MasterDB() *D {
	return any(r.dbMaster).(*D)
}

func (r *baseRepo[D, I, E]) SlaveDB() *D {
	return any(r.dbSlave).(*D)
}

func (r *baseRepo[D, I, E]) FindAll(ctx context.Context, filter map[string]any) (res []E, err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"filter": filter,
	}).AddAttribute("table.name", r.Entity.TableName())

	defer func() {
		span.SetFunctionOutput(tracer.Metadata{
			"result": res,
		}).End(err)
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

func (r *baseRepo[D, I, E]) FindById(ctx context.Context, ID I) (res *E, err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"id": ID,
	}).AddAttribute("table.name", r.Entity.TableName())

	defer func() {
		span.SetFunctionOutput(tracer.Metadata{
			"result": res,
		}).End(err)
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

func (r *baseRepo[D, I, E]) FindByIdAndLock(ctx context.Context, ID I, trx *D) (res *E, err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"id": ID,
	}).AddAttribute("table.name", r.Entity.TableName())

	defer func() {
		span.SetFunctionOutput(tracer.Metadata{
			"result": res,
		}).End(err)
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

func (r *baseRepo[D, I, E]) FindByIds(ctx context.Context, IDs []I) (res []E, err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"ids": IDs,
	}).AddAttribute("table.name", r.Entity.TableName())

	defer func() {
		span.SetFunctionOutput(tracer.Metadata{
			"result": res,
		}).End(err)
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

func (r *baseRepo[D, I, E]) FindByOffset(ctx context.Context, filter map[string]any, sort []string, size int, page int) (res database.Pagination[E], err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"filter": filter,
		"sort":   sort,
		"size":   size,
		"page":   page,
	}).AddAttribute("table.name", r.Entity.TableName())

	defer func() {
		span.SetFunctionOutput(tracer.Metadata{
			"result": res,
		}).End(err)
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

func (r *baseRepo[D, I, E]) FindByCursor(ctx context.Context, filter map[string]any, sort []string, size int, next *I) (res database.Pagination[E], err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"filter": filter,
		"sort":   sort,
		"size":   size,
		"next":   next,
	}).AddAttribute("table.name", r.Entity.TableName())

	defer func() {
		span.SetFunctionOutput(tracer.Metadata{
			"result": res,
		}).End(err)
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
func (r *baseRepo[D, I, E]) Insert(ctx context.Context, payload E, trx *D) (res E, err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"payload": payload,
	}).AddAttribute("table.name", r.Entity.TableName())

	defer func() {
		span.SetFunctionOutput(tracer.Metadata{
			"result": res,
		}).End(err)
	}()

	db := r.dbMaster
	if trx, ok := any(trx).(*gorm.DB); ok {
		db = trx
	}

	err = db.WithContext(ctx).Create(&payload).Error
	if err != nil {
		return payload, err
	}

	return payload, nil
}

// TODO: Check 'res' is still necessary
func (r *baseRepo[D, I, E]) InsertMany(ctx context.Context, payload []E, trx *D) (res []E, err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"payload": payload,
	}).AddAttribute("table.name", r.Entity.TableName())

	defer func() {
		span.SetFunctionOutput(tracer.Metadata{
			"result": res,
		}).End(err)
	}()

	db := r.dbMaster
	if trx, ok := any(trx).(*gorm.DB); ok {
		db = trx
	}

	err = db.WithContext(ctx).CreateInBatches(payload, config.MySQL.InsertBatchSize).Error
	if err != nil {
		return payload, err
	}

	return payload, nil
}

func (r *baseRepo[D, I, E]) Update(ctx context.Context, payload E, trx *D) (err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"payload": payload,
	}).AddAttribute("table.name", r.Entity.TableName())

	defer func() {
		span.End(err)
	}()

	db := r.dbMaster
	if trx, ok := any(trx).(*gorm.DB); ok {
		db = trx
	}

	err = db.WithContext(ctx).Save(&payload).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *baseRepo[D, I, E]) UpdateById(ctx context.Context, ID I, payload map[string]any, trx *D) (res E, err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"id":      ID,
		"payload": payload,
	}).AddAttribute("table.name", r.Entity.TableName())

	defer func() {
		span.SetFunctionOutput(tracer.Metadata{
			"result": res,
		}).End(err)
	}()

	db := r.dbMaster
	if trx, ok := any(trx).(*gorm.DB); ok {
		db = trx
	}

	err = db.WithContext(ctx).Model(&res).Where("id=?", ID).Updates(payload).Scan(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}

// TODO: Check 'res' is needed
func (r *baseRepo[D, I, E]) UpdateByIds(ctx context.Context, IDs []I, payload map[string]any, trx *D) (err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"ids":     IDs,
		"payload": payload,
	}).AddAttribute("table.name", r.Entity.TableName())

	defer func() {
		span.End(err)
	}()

	db := r.dbMaster
	if trx, ok := any(trx).(*gorm.DB); ok {
		db = trx
	}

	err = db.WithContext(ctx).Model(&r.Entity).Where("id IN ?", IDs).Updates(payload).Error
	if err != nil {
		return err
	}

	return nil
}

// TODO: Check 'res' is needed
func (r *baseRepo[D, I, E]) UpdateMany(ctx context.Context, filter map[string]any, payload map[string]any, trx *D) (err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"filter":  filter,
		"payload": payload,
	}).AddAttribute("table.name", r.Entity.TableName())

	defer func() {
		span.End(err)
	}()

	db := r.dbMaster
	if trx, ok := any(trx).(*gorm.DB); ok {
		db = trx
	}

	err = db.WithContext(ctx).Model(&r.Entity).Where(filter).Updates(payload).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *baseRepo[D, I, E]) DeleteById(ctx context.Context, ID I, trx *D) (err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"id": ID,
	}).AddAttribute("table.name", r.Entity.TableName())

	defer func() {
		span.End(err)
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

func (r *baseRepo[D, I, E]) DeleteByIds(ctx context.Context, IDs []I, trx *D) (err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"ids": IDs,
	}).AddAttribute("table.name", r.Entity.TableName())

	defer func() {
		span.End(err)
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

func (r *baseRepo[D, I, E]) DeleteMany(ctx context.Context, filter map[string]any, trx *D) (err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"filter": filter,
	}).AddAttribute("table.name", r.Entity.TableName())

	defer func() {
		span.End(err)
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

func (r *baseRepo[D, I, E]) Begin(ctx context.Context) (trx *D, err error) {
	db := r.dbMaster.WithContext(ctx).Begin()
	if db.Error != nil {
		return trx, db.Error
	}

	return any(db).(*D), nil
}

func (r *baseRepo[D, I, E]) Rollback(trx *D) *D {
	db, ok := any(trx).(*gorm.DB)
	if !ok {
		return trx
	}

	db = db.Rollback()

	return any(db).(*D)
}

func (r *baseRepo[D, I, E]) Commit(trx *D) *D {
	db, ok := any(trx).(*gorm.DB)
	if !ok {
		return trx
	}

	db = db.Commit()

	return any(db).(*D)
}
