package services

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/ms-accounts-grpc/fn/fnFakeit"
	"github.com/d3v-friends/ms-accounts-grpc/models"
	"github.com/d3v-friends/ms-accounts-grpc/vars"
	"github.com/d3v-friends/pure-go/fnEnv"
	"github.com/d3v-friends/pure-go/fnLogger"
	"github.com/d3v-friends/pure-go/fnPanic"
	"github.com/d3v-friends/pure-go/fnParams"
	"go.mongodb.org/mongo-driver/mongo"
)

type testTools struct {
	DB *mongo.Database
}

func newTestTools(isInit ...bool) (res *testTools) {
	var ctx = context.TODO()
	fnFakeit.Init()
	res = &testTools{}

	var mClient = fnPanic.OnValue(mango.NewClient(&mango.ClientOpt{
		Host:     fnEnv.Read("MG_HOST"),
		Username: fnEnv.Read("MG_USERNAME"),
		Password: fnEnv.Read("MG_PASSWORD"),
		Database: fnEnv.Read("MG_DATABASE"),
	}))

	res.DB = mClient.Database()

	if fnParams.Get(isInit) {
		fnPanic.On(res.DB.Drop(ctx))
		fnPanic.On(mClient.Migrate(ctx, models.All...))
		fnPanic.On(res.indexAccount())
	}

	return
}

func (x *testTools) context() context.Context {
	return vars.SetUtils(context.TODO(), x.DB, fnLogger.NewDefaultLogger())
}

func (x *testTools) createAccount() (res *ICreateAccount) {
	res = &ICreateAccount{
		Identifier: map[string]string{
			"email": gofakeit.Email(),
		},
		Property: map[string]string{
			"address": gofakeit.Address().Address,
		},
		Verifier: map[string]*IVerifier{
			"passwd": {
				Salt:   "passwd",
				Passwd: "passwd",
				Etc:    "passwd",
				Mode:   VerifyMode_COMPARE,
			},
		},
		Permission: map[string]bool{
			"signIn": true,
		},
	}

	return
}

func (x *testTools) indexAccount() (err error) {
	var sv = &SystemImpl{}
	_, err = sv.UpdateKeys(x.context(), &IUpdateKeys{
		Identifier: []string{"email"},
		Property:   []string{"address"},
		Permission: []string{"signIn"},
	})
	return
}

func (x *testTools) iCreateSession(account *Account) (res *ICreateSession) {
	res = &ICreateSession{
		Identifier: map[string]string{
			"email": account.Identifier["email"],
		},
		Verifier: map[string]string{
			"passwd": "passwd",
		},
		Permission: map[string]bool{
			"signIn": true,
		},
	}
	return
}
