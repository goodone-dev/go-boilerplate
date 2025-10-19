package mongodb

import (
	"context"
	"errors"

	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type BaseRepo[D any, I any, E database.Entity] struct {
	Entity   E
	dbMaster *mongo.Database
	dbSlave  *mongo.Database
}

func NewBaseRepo[D any, I any, E database.Entity](dbConn mongoConnection) database.IBaseRepository[D, I, E] {
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

	coll := r.dbSlave.Collection(r.Entity.TableName())

	filter["deleted_at"] = nil
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &res)
	return
}

func (r *BaseRepo[D, I, E]) FindById(ctx context.Context, ID I) (res *E, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, ID)
	defer func() {
		span.EndSpan(err, res)
	}()

	coll := r.dbSlave.Collection(r.Entity.TableName())

	err = coll.FindOne(ctx, bson.M{"_id": ID, "deleted_at": nil}).Decode(&res)
	return
}

func (r *BaseRepo[D, I, E]) FindByIdAndLock(ctx context.Context, ID I, trx *D) (res *E, err error) {
	return nil, errors.New("locking not supported")
}

func (r *BaseRepo[D, I, E]) FindByIds(ctx context.Context, IDs []I) (res []E, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, IDs)
	defer func() {
		span.EndSpan(err, res)
	}()

	coll := r.dbSlave.Collection(r.Entity.TableName())

	cursor, err := coll.Find(ctx, bson.M{"_id": bson.M{"$in": IDs}, "deleted_at": nil})
	if err != nil {
		return
	}

	err = cursor.All(ctx, &res)
	return
}

func (r *BaseRepo[D, I, E]) OffsetPagination(ctx context.Context, filter map[string]any, sort []string, page int, size int) (res database.Pagination[E], err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, filter, sort, page, size)
	defer func() {
		span.EndSpan(err, res)
	}()

	coll := r.dbSlave.Collection(r.Entity.TableName())

	count, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return
	}

	opt := options.Find().
		SetSort(sort).
		SetLimit(int64(size)).
		SetSkip(int64((page - 1) * size))

	filter["deleted_at"] = nil

	cursor, err := coll.Find(ctx, filter, opt)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &res.Data)
	if err != nil {
		return
	}

	res.Metadata.Total = int(count)
	res.Metadata.Pages = (int(count) + size - 1) / size
	res.Metadata.Page = page
	res.Metadata.Size = size

	return
}

func (r *BaseRepo[D, I, E]) Insert(ctx context.Context, model E, trx *D) (res E, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, model)
	defer func() {
		span.EndSpan(err, res)
	}()

	coll := r.dbMaster.Collection(r.Entity.TableName())

	result, err := coll.InsertOne(ctx, model)
	if err != nil {
		return
	}

	err = coll.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&res)
	return
}

func (r *BaseRepo[D, I, E]) InsertMany(ctx context.Context, models []E, trx *D) (res []E, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, models)
	defer func() {
		span.EndSpan(err, res)
	}()

	coll := r.dbMaster.Collection(r.Entity.TableName())

	docs := make([]any, 0, len(models))
	for _, model := range models {
		docs = append(docs, model)
	}

	result, err := coll.InsertMany(ctx, docs)

	cursor, err := coll.Find(ctx, bson.M{"_id": bson.M{"$in": result.InsertedIDs}})
	if err != nil {
		return
	}

	err = cursor.All(ctx, &res)
	return
}

func (r *BaseRepo[D, I, E]) UpdateById(ctx context.Context, ID I, payload map[string]any, trx *D) (res E, err error) {
	_, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, payload)
	defer func() {
		span.EndSpan(err, res)
	}()

	coll := r.dbMaster.Collection(r.Entity.TableName())

	err = coll.FindOneAndUpdate(ctx, bson.M{"_id": ID}, payload).Decode(&res)
	return
}

func (r *BaseRepo[D, I, E]) UpdateByIds(ctx context.Context, IDs []I, payload map[string]any, trx *D) (err error) {
	_, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, payload)
	defer func() {
		span.EndSpan(err)
	}()

	coll := r.dbMaster.Collection(r.Entity.TableName())

	_, err = coll.UpdateMany(ctx, bson.M{"_id": IDs}, payload)
	return
}

func (r *BaseRepo[D, I, E]) DeleteById(ctx context.Context, ID I, trx *D) (err error) {
	_, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, ID)
	defer func() {
		span.EndSpan(err)
	}()

	coll := r.dbMaster.Collection(r.Entity.TableName())

	_, err = coll.DeleteOne(ctx, bson.M{"_id": ID})
	return
}

func (r *BaseRepo[D, I, E]) DeleteByIds(ctx context.Context, IDs []I, trx *D) (err error) {
	_, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, IDs)
	defer func() {
		span.EndSpan(err)
	}()

	coll := r.dbMaster.Collection(r.Entity.TableName())

	_, err = coll.DeleteOne(ctx, bson.M{"_id": IDs})
	return
}

func (r *BaseRepo[D, I, E]) Begin(ctx context.Context) (*D, error) {
	return nil, errors.New("transaction not supported")
}

func (r *BaseRepo[D, I, E]) Rollback(trx *D) *D {
	return nil
}

func (r *BaseRepo[D, I, E]) Commit(trx *D) *D {
	return nil
}
