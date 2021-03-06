package storage

import (
	"context"

	"github.com/huzhongqing/qelog/pkg/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 当报警规则超过上千的规则，可优化语句
func (store *Store) FindAllEnableAlarmRule(ctx context.Context) ([]*model.AlarmRule, error) {
	docs := make([]*model.AlarmRule, 0)
	coll := store.database.Collection(model.CollectionNameAlarmRule)
	err := store.database.Find(ctx, coll, bson.M{"enable": true}, &docs)
	return docs, handlerError(err)
}

func (store *Store) FindAlarmRuleList(ctx context.Context, filter bson.M, result interface{}, opt *options.FindOptions) (int64, error) {
	c, err := store.database.FindAndCount(ctx, store.database.Collection(model.CollectionNameAlarmRule), filter, result, opt)
	return c, handlerError(err)
}

func (store *Store) InsertAlarmRule(ctx context.Context, doc *model.AlarmRule) error {
	_, err := store.database.Collection(doc.CollectionName()).InsertOne(ctx, doc)
	return handlerError(err)
}

func (store *Store) FindOneAlarmRule(ctx context.Context, filter bson.M, doc *model.AlarmRule) (bool, error) {
	return store.database.FindOne(ctx, store.database.Collection(doc.CollectionName()), filter, doc)
}

func (store *Store) UpdateAlarmRule(ctx context.Context, filter, update bson.M) error {
	uRet, err := store.database.Collection(model.CollectionNameAlarmRule).UpdateOne(ctx, filter, update)
	if err != nil {
		return handlerError(err)
	}
	if uRet.MatchedCount <= 0 {
		return ErrNotMatched
	}
	return nil
}

func (store *Store) UpdateManyAlarmRule(ctx context.Context, filter, update bson.M) error {
	_, err := store.database.Collection(model.CollectionNameAlarmRule).UpdateMany(ctx, filter, update)
	if err != nil {
		return handlerError(err)
	}
	return nil
}

func (store *Store) DeleteAlarmRule(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{
		"_id": id,
	}
	_, err := store.database.Collection(model.CollectionNameAlarmRule).DeleteOne(ctx, filter)
	return handlerError(err)
}

func (store *Store) FindHookURL(ctx context.Context, filter bson.M, result interface{}, opt *options.FindOptions) (int64, error) {
	c, err := store.database.FindAndCount(ctx, store.database.Collection(model.CollectionNameHookURL), filter, result, opt)
	return c, handlerError(err)
}

func (store *Store) FindOneHookURL(ctx context.Context, filter bson.M, doc *model.HookURL) (bool, error) {
	return store.database.FindOne(ctx, store.database.Collection(doc.CollectionName()), filter, doc)
}

func (store *Store) InsertHookURL(ctx context.Context, doc *model.HookURL) error {
	_, err := store.database.Collection(doc.CollectionName()).InsertOne(ctx, doc)
	return handlerError(err)
}

func (store *Store) UpdateHookURL(ctx context.Context, filter, update bson.M) error {
	uRet, err := store.database.Collection(model.CollectionNameHookURL).UpdateOne(ctx, filter, update)
	if err != nil {
		return handlerError(err)
	}
	if uRet.MatchedCount <= 0 {
		return ErrNotMatched
	}
	return nil
}

func (store *Store) DelHookURL(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{
		"_id": id,
	}
	_, err := store.database.Collection(model.CollectionNameHookURL).DeleteOne(ctx, filter)
	return handlerError(err)
}
