package models

import (
	"context"
	"github.com/d3v-friends/mango/mtype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type (
	System struct {
		Id      primitive.ObjectID `bson:"_id"`
		Data    SystemData         `bson:"data"`
		History []*SystemData      `bson:"history"`
	}

	SystemData struct {
		Identifier []string  `bson:"identifier"`
		Property   []string  `bson:"property"`
		Permission []string  `bson:"permission"`
		CreatedAt  time.Time `bson:"createdAt"`
	}
)

const (
	colSystem = "systems"

	fSystemId             = "_id"
	fSystemData           = "data"
	fSystemHistory        = "history"
	fSystemDataIdentifier = "data.identifier"
	fSystemDataProperty   = "data.property"
	fSystemDataPermission = "data.permission"
	fSystemDataCreatedAt  = "data.createdAt"
)

var mgSystem = mtype.FnMigrateList{
	func(ctx context.Context, collection *mongo.Collection) (migrationNm string, err error) {
		var now = time.Now()
		migrationNm = "init indexing"
		_, err = collection.InsertOne(ctx, &System{
			Id: primitive.NilObjectID,
			Data: SystemData{
				Identifier: []string{},
				Property:   []string{},
				Permission: []string{},
				CreatedAt:  now,
			},
			History: make([]*SystemData, 0),
		})
		return
	},
}

func (x System) GetCollectionNm() string {
	return colSystem
}

func (x System) GetMigrateList() mtype.FnMigrateList {
	return mgSystem
}
