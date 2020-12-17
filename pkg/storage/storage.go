package storage

import (
	"context"
	"fmt"

	"github.com/huzhongqing/qelog/libs/mongo"
)

type Store struct {
	database *mongo.Database
}

func New(database *mongo.Database) *Store {
	store := &Store{
		database: database,
	}
	return store
}
func (store *Store) ListAllCollectionNames(ctx context.Context) ([]string, error) {
	names, err := store.database.ListAllCollectionNames(ctx)
	return names, WrapError(err)
}

func (store *Store) UpsertCollectionIndexMany(indexs []mongo.Index) error {
	err := store.database.UpsertCollectionIndexMany(indexs)
	return WrapError(err)
}

func WrapError(err error) error {
	if err != nil {
		fmt.Println(err)
	}
	return err
}
