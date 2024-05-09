package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/akbariandev/pacassistant/internal/domain"
	"github.com/akbariandev/pacassistant/pkg/logger"
	"github.com/akbariandev/pacassistant/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository[T any] struct {
	logger     logger.Logger
	collection *mongo.Collection
}

func newRepository[T any](
	collection *mongo.Collection,
	logger logger.Logger,
) domain.Repository[T] {
	return &Repository[T]{
		collection: collection,
		logger:     logger,
	}
}

func (r *Repository[T]) Get(ctx context.Context, id primitive.ObjectID) (T, error) {
	var model T
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&model); err != nil {
		return model, err
	}

	return model, nil
}

func (r *Repository[T]) Exists(ctx context.Context, filter interface{}) bool {
	var obj interface{}
	err := r.collection.FindOne(ctx, filter).Decode(&obj)
	if errors.Is(err, mongo.ErrNoDocuments) || err != nil {
		return false
	}
	return true
}

func (r *Repository[T]) Create(ctx context.Context, model *T) error {
	if _, err := r.collection.InsertOne(ctx, model); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("already exists")
		}
		r.logger.ErrorContext(ctx, true, err.Error())

		return errors.New("database error happen")
	}

	return nil
}

func (r *Repository[T]) CreateMany(ctx context.Context, models []any) error {
	if _, err := r.collection.InsertMany(ctx, models); err != nil {
		return err
	}
	return nil
}

func (r *Repository[T]) Replace(ctx context.Context, model *T, id primitive.ObjectID) error {
	if _, err := r.collection.ReplaceOne(ctx, bson.M{"_id": id}, model); err != nil {
		return err
	}
	return nil
}

func (r *Repository[T]) InsertOrReplace(ctx context.Context, model *T, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": model}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[T]) UpdateFields(ctx context.Context, model *T, id primitive.ObjectID, fields ...string) error {
	return mongodb.UpdateFields(ctx, r.collection, bson.M{"_id": id}, model, fields...)
}

func (r *Repository[T]) Delete(ctx context.Context, id primitive.ObjectID) error {
	if _, err := r.collection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return err
	}

	return nil
}
