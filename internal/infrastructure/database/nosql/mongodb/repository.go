package mongodb

import (
	"context"

	database "github.com/BagusAK95/go-boilerplate/internal/infrastructure/database/nosql"
	"github.com/BagusAK95/go-boilerplate/internal/utils/tracer"
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseRepo[M database.Entity] struct {
	Entity   M
	dbMaster *mongo.Database
	dbSlave  *mongo.Database
}

func NewBaseRepo[M database.Entity](dbConn dbConnection) database.IBaseRepository[M] {
	return &BaseRepo[M]{
		dbMaster: dbConn.Master,
		dbSlave:  dbConn.Slave,
	}
}

func (r *BaseRepo[M]) InsertOne(ctx context.Context, doc M) (res *database.InsertOneResult, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, doc)
	defer func() {
		span.EndSpan(err, res)
	}()

	collection := r.dbMaster.Collection(r.Entity.CollectionName())

	rest, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return
	}

	return (*database.InsertOneResult)(rest), nil
}

func (r *BaseRepo[M]) InsertMany(ctx context.Context, models []M) (res *database.InsertManyResult, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, models)
	defer func() {
		span.EndSpan(err, res)
	}()

	docs := make([]any, 0, len(models))
	for _, model := range models {
		docs = append(docs, model)
	}

	collection := r.dbMaster.Collection(r.Entity.CollectionName())

	rest, err := collection.InsertMany(ctx, docs)

	return (*database.InsertManyResult)(rest), nil
}

func (r *BaseRepo[M]) FindOne(ctx context.Context, filter any) (res M, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, filter)
	defer func() {
		span.EndSpan(err, res)
	}()

	collection := r.dbSlave.Collection(r.Entity.CollectionName())

	err = collection.FindOne(ctx, filter).Decode(&res)

	return
}

func (r *BaseRepo[M]) Find(ctx context.Context, filter any) (res []M, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, filter)
	defer func() {
		span.EndSpan(err, res)
	}()

	collection := r.dbSlave.Collection(r.Entity.CollectionName())

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &res)

	return
}

func (r *BaseRepo[M]) UpdateOne(ctx context.Context, filter any, data any) (res *database.UpdateResult, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, filter)
	defer func() {
		span.EndSpan(err, res)
	}()

	collection := r.dbMaster.Collection(r.Entity.CollectionName())

	rest, err := collection.UpdateOne(ctx, filter, data)
	if err != nil {
		return
	}

	return (*database.UpdateResult)(rest), nil
}

func (r *BaseRepo[M]) UpdateMany(ctx context.Context, filter any, data any) (res *database.UpdateResult, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, filter)
	defer func() {
		span.EndSpan(err, res)
	}()

	collection := r.dbMaster.Collection(r.Entity.CollectionName())

	rest, err := collection.UpdateMany(ctx, filter, data)
	if err != nil {
		return
	}

	return (*database.UpdateResult)(rest), nil
}

func (r *BaseRepo[M]) DeleteOne(ctx context.Context, filter any) (res *database.DeleteResult, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, filter)
	defer func() {
		span.EndSpan(err, res)
	}()

	collection := r.dbMaster.Collection(r.Entity.CollectionName())

	rest, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return
	}

	return (*database.DeleteResult)(rest), nil
}

func (r *BaseRepo[M]) DeleteMany(ctx context.Context, filter any) (res *database.DeleteResult, err error) {
	ctx, span := tracer.SpanPrefixName(r.Entity.RepositoryName()).StartSpan(ctx, filter)
	defer func() {
		span.EndSpan(err, res)
	}()

	collection := r.dbMaster.Collection(r.Entity.CollectionName())

	rest, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return
	}

	return (*database.DeleteResult)(rest), nil
}
