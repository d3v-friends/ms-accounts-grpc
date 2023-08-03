package models

import (
	"context"
	"fmt"
	"github.com/d3v-friends/mango/mvars"
	"github.com/d3v-friends/ms-accounts-grpc/vars"
	"github.com/d3v-friends/pure-go/fnReflect"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type (
	IfUpdateSystem interface {
		GetIdentifier() []string
		GetProperty() []string
		GetPermission() []string
	}
)

func FindSystem(ctx context.Context) (system *System, err error) {
	db := vars.GetDB(ctx)
	system = &System{}

	resp := db.
		Collection(system.GetCollectionNm()).
		FindOne(ctx, &bson.M{
			mvars.FID: primitive.NilObjectID,
		})

	if err = resp.Err(); err != nil {
		return
	}

	if err = resp.Decode(system); err != nil {
		return
	}

	return
}

func UpdateSystem(ctx context.Context, i IfUpdateSystem) (system *System, err error) {
	db := vars.GetDB(ctx)

	resp := db.Collection(colSystem).FindOneAndUpdate(ctx, &bson.M{}, &bson.M{
		mvars.OSet: bson.M{
			fSystemDataIdentifier: i.GetIdentifier(),
			fSystemDataProperty:   i.GetProperty(),
			fSystemDataPermission: i.GetPermission(),
			fSystemDataCreatedAt:  time.Now(),
		},
	})

	if err = resp.Err(); err != nil {
		return
	}

	system = &System{}
	if err = resp.Decode(system); err != nil {
		return
	}

	if _, err = db.Collection(colSystem).UpdateOne(ctx, &bson.M{}, &bson.M{
		mvars.OPush: bson.M{
			fSystemHistory: system.Data,
		},
	}); err != nil {
		return
	}

	// todo 나중에 여기 다시 한번 확인해보기
	// accounts collection reindexing
	_, _ = db.
		Collection(colAccount).
		Indexes().
		DropAll(ctx)

	idxList := make([]mongo.IndexModel, 0)

	// identifier
	idxList = append(idxList, createIndex(fAccount2DataIdentifier, i.GetIdentifier(), &options.IndexOptions{
		Unique: fnReflect.ToPointer(true),
	})...)

	// property
	idxList = append(idxList, createIndex(fAccount2DataProperty, i.GetProperty(), nil)...)

	// permission
	idxList = append(idxList, createIndex(fAccount2DataPermission, i.GetPermission(), nil)...)

	// session
	idxList = append(idxList, mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   fAccountSessionId,
				Value: 1,
			},
		},
		Options: &options.IndexOptions{
			Unique: fnReflect.ToPointer(true),
		},
	})

	if _, err = db.
		Collection(colAccount).
		Indexes().
		CreateMany(ctx, idxList); err != nil {
		return
	}

	return FindSystem(ctx)
}

func createIndex(prefix string, keyList []string, options *options.IndexOptions) (list []mongo.IndexModel) {
	if len(keyList) == 0 {
		return
	}

	for _, key := range keyList {
		list = append(list, mongo.IndexModel{
			Keys: &bson.D{
				{fmt.Sprintf("%s.%s", prefix, key), 1},
			},
			Options: options,
		})
	}
	return
}
