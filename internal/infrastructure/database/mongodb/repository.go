package mongodb

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"time"

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

func (r *BaseRepo[D, I, E]) FindAll(ctx context.Context, filter map[string]any) (res []E, err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, filter)
	defer func() {
		span.Stop(err, res)
	}()

	coll := r.dbSlave.Collection(r.Entity.TableName())

	filter["deleted_at"] = nil

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &res)
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

	coll := r.dbSlave.Collection(r.Entity.TableName())

	err = coll.FindOne(ctx, bson.M{"_id": ID, "deleted_at": nil}).Decode(&res)
	if err != nil {
		return
	}

	return
}

func (r *BaseRepo[D, I, E]) FindByIdAndLock(ctx context.Context, ID I, trx *D) (res *E, err error) {
	return nil, errors.New("locking not supported")
}

func (r *BaseRepo[D, I, E]) FindByIds(ctx context.Context, IDs []I) (res []E, err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, IDs)
	defer func() {
		span.Stop(err, res)
	}()

	coll := r.dbSlave.Collection(r.Entity.TableName())

	cursor, err := coll.Find(ctx, bson.M{"_id": bson.M{"$in": IDs}, "deleted_at": nil})
	if err != nil {
		return
	}

	err = cursor.All(ctx, &res)
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

	coll := r.dbSlave.Collection(r.Entity.TableName())

	filter["deleted_at"] = nil
	count, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return
	}

	if size <= 0 {
		size = 10
	}
	if page <= 0 {
		page = 1
	}

	opt := options.Find().
		SetSort(sort).
		SetLimit(int64(size)).
		SetSkip(int64((page - 1) * size))

	cursor, err := coll.Find(ctx, filter, opt)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &res.Data)
	if err != nil {
		return
	}

	var pages int
	if count > 0 {
		pages = int(math.Ceil(float64(count) / float64(size)))
	}

	res.Metadata.Total = &count
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

	coll := r.dbSlave.Collection(r.Entity.TableName())

	filter["deleted_at"] = nil
	if next != nil {
		filter["id > ?"] = *next
	}

	count, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return
	}

	if size <= 0 {
		size = 10
	}

	opt := options.Find().
		SetSort(sort).
		SetLimit(int64(size))

	cursor, err := coll.Find(ctx, filter, opt)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &res.Data)
	if err != nil {
		return
	}

	var pages int
	if count > 0 {
		pages = int(math.Ceil(float64(count) / float64(size)))
	}

	res.Metadata.Total = &count
	res.Metadata.Pages = &pages
	res.Metadata.Size = &size

	return
}

func (r *BaseRepo[D, I, E]) Insert(ctx context.Context, model E, trx *D) (res E, err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, model)
	defer func() {
		span.Stop(err, res)
	}()

	coll := r.dbMaster.Collection(r.Entity.TableName())

	result, err := coll.InsertOne(ctx, model)
	if err != nil {
		return
	}

	err = coll.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&res)
	if err != nil {
		return
	}

	return
}

func (r *BaseRepo[D, I, E]) InsertMany(ctx context.Context, models []E, trx *D) (res []E, err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, models)
	defer func() {
		span.Stop(err, res)
	}()

	coll := r.dbMaster.Collection(r.Entity.TableName())

	result, err := coll.InsertMany(ctx, models)
	if err != nil {
		return
	}

	cursor, err := coll.Find(ctx, bson.M{"_id": bson.M{"$in": result.InsertedIDs}})
	if err != nil {
		return
	}

	err = cursor.All(ctx, &res)
	if err != nil {
		return
	}

	return
}

func (r *BaseRepo[D, I, E]) Update(ctx context.Context, model E, trx *D) (err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, model)
	defer func() {
		span.Stop(err)
	}()

	coll := r.dbMaster.Collection(r.Entity.TableName())

	data, err := json.Marshal(model)
	if err != nil {
		return
	}

	var req database.BaseEntity[I]
	err = json.Unmarshal(data, &req)
	if err != nil {
		return
	}

	_, err = coll.UpdateOne(ctx, bson.M{"_id": req.ID}, bson.M{"$set": model})
	if err != nil {
		return
	}

	return
}

func (r *BaseRepo[D, I, E]) UpdateById(ctx context.Context, ID I, payload map[string]any, trx *D) (res E, err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, payload)
	defer func() {
		span.Stop(err, res)
	}()

	coll := r.dbMaster.Collection(r.Entity.TableName())

	err = coll.FindOneAndUpdate(ctx, bson.M{"_id": ID}, bson.M{"$set": payload}).Decode(&res)
	if err != nil {
		return
	}

	return
}

func (r *BaseRepo[D, I, E]) UpdateByIds(ctx context.Context, IDs []I, payload map[string]any, trx *D) (err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, payload)
	defer func() {
		span.Stop(err)
	}()

	coll := r.dbMaster.Collection(r.Entity.TableName())

	_, err = coll.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": IDs}}, bson.M{"$set": payload})
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepo[D, I, E]) UpdateMany(ctx context.Context, filter map[string]any, payload map[string]any, trx *D) (err error) {
	ctx, span := tracer.PrefixName(r.Entity.RepositoryName()).Start(ctx, filter, payload)
	defer func() {
		span.Stop(err)
	}()

	coll := r.dbMaster.Collection(r.Entity.TableName())

	_, err = coll.UpdateMany(ctx, filter, bson.M{"$set": payload})
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

	coll := r.dbMaster.Collection(r.Entity.TableName())

	_, err = coll.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
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

	coll := r.dbMaster.Collection(r.Entity.TableName())

	_, err = coll.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": IDs}}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
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

	coll := r.dbMaster.Collection(r.Entity.TableName())

	_, err = coll.UpdateMany(ctx, filter, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	if err != nil {
		return err
	}

	return nil
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
