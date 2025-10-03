package nosql

import (
	"context"
)

type IBaseRepository[M Entity] interface {
	InsertOne(ctx context.Context, req M) (res *InsertOneResult, err error)
	InsertMany(ctx context.Context, models []M) (res *InsertManyResult, err error)
	FindOne(ctx context.Context, filter any) (res M, err error)
	Find(ctx context.Context, filter any) (res []M, err error)
	UpdateOne(ctx context.Context, filter any, data any) (res *UpdateResult, err error)
	UpdateMany(ctx context.Context, filter any, data any) (res *UpdateResult, err error)
	DeleteOne(ctx context.Context, filter any) (res *DeleteResult, err error)
	DeleteMany(ctx context.Context, filter any) (res *DeleteResult, err error)
}

type InsertOneResult struct {
	InsertedID any
}

type InsertManyResult struct {
	InsertedIDs []any
}

type UpdateResult struct {
	MatchedCount  int64
	ModifiedCount int64
	UpsertedCount int64
	UpsertedID    any
}

type DeleteResult struct {
	DeletedCount int64 `bson:"n"`
}
