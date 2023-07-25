package services

import (
	"context"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/ms-accounts-grpc/fn/fnFakeit"
	"github.com/d3v-friends/ms-accounts-grpc/models"
	"github.com/d3v-friends/ms-accounts-grpc/vars"
	"github.com/d3v-friends/pure-go/fnEnv"
	"github.com/d3v-friends/pure-go/fnLogger"
	"github.com/d3v-friends/pure-go/fnPanic"
	"go.mongodb.org/mongo-driver/mongo"
)

type testTools struct {
	DB *mongo.Database
}

func newTestTools() (res *testTools) {
	fnFakeit.Init()
	res = &testTools{}

	var mClient = fnPanic.OnValue(mango.NewClient(&mango.ClientOpt{
		Host:     fnEnv.Read("MG_HOST"),
		Username: fnEnv.Read("MG_USERNAME"),
		Password: fnEnv.Read("MG_PASSWORD"),
		Database: fnEnv.Read("MG_DATABASE"),
	}))

	res.DB = mClient.Database()

	fnPanic.On(mClient.Migrate(context.TODO(), models.All...))

	return
}

func (x *testTools) Context() context.Context {
	return vars.SetUtils(context.TODO(), x.DB, fnLogger.NewDefaultLogger())
}
