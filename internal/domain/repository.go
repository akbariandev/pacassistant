package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository[T any] interface {
	Get(ctx context.Context, id primitive.ObjectID) (T, error)
	Exists(ctx context.Context, filter interface{}) bool
	Create(ctx context.Context, model *T) error
	CreateMany(ctx context.Context, models []any) error
	Replace(ctx context.Context, model *T, id primitive.ObjectID) error
	InsertOrReplace(ctx context.Context, model *T, id primitive.ObjectID) error
	UpdateFields(ctx context.Context, model *T, id primitive.ObjectID, fields ...string) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
