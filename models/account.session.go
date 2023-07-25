package models

import (
	"context"
	"fmt"
	"github.com/d3v-friends/mango/mvars"
	"github.com/d3v-friends/ms-accounts-grpc/fn/fnOTP"
	"github.com/d3v-friends/ms-accounts-grpc/vars"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

type ICreateSession struct {
	Identifier map[string]string
	Permission map[string]bool
	Verifier   map[string]string
}

func CreateSession(ctx context.Context, i *ICreateSession) (res *SessionData, err error) {
	var db, logger = vars.GetUtils(ctx)
	var now = time.Now()

	var account *Account
	if account, err = ReadAccountOne(ctx, &IReadAccount{
		Identifier: i.Identifier,
		Permission: i.Permission,
	}); err != nil {
		return
	}

	for key, passwd := range i.Verifier {
		verifier, has := account.Account.Data.Verifier[key]
		if !has {
			err = fmt.Errorf(
				"not found verifier: accountId=%s, key=%s",
				account.Id.Hex(),
				key,
			)
			return
		}

		switch strings.ToLower(verifier.Mode) {
		case "compare":
			if verifier.Passwd != passwd {
				err = fmt.Errorf(
					"invalid password: accountId=%s, key=%s",
					account.Id.Hex(),
					key,
				)
				return
			}
		case "g_otp":
			if !fnOTP.Verify(verifier.Salt, passwd) {
				err = fmt.Errorf(
					"invalid otp: otp=%s",
					passwd,
				)
				return
			}
		default:
			err = fmt.Errorf(
				"invalid verifier mode: mode=%s",
				verifier.Mode,
			)
			return
		}
	}

	var sessionId = primitive.NewObjectID()
	if _, err = db.Collection(colAccount).UpdateOne(ctx,
		bson.M{
			fAccountId: account.Id,
		},
		bson.M{
			mvars.OPush: bson.M{
				fAccountSessionData: sessionId,
				fAccountSessionHistories: SessionData{
					Id:        sessionId,
					CreatedAt: now,
				},
			},
		}); err != nil {
		return
	}

	logger.Trace("created session: sessionId=%s", sessionId.Hex())

	res = &SessionData{
		Id:        sessionId,
		CreatedAt: now,
	}

	return
}

func VerifySession(ctx context.Context, i *IVerifyAccountSession) (account *Account, err error) {
	var _, logger = vars.GetUtils(ctx)
	if account, err = ReadAccountOne(ctx, &IReadAccount{
		Permission: i.Permission,
		SessionId:  &i.SessionId,
	}); err != nil {
		return
	}

	logger.Trace(
		"success verify account: accountId=%s, sessionId=%s",
		account.Id.Hex(),
		i.SessionId.Hex(),
	)

	return
}

type IDeleteSession struct {
	SessionId primitive.ObjectID
	AccountId primitive.ObjectID
}

func (x IDeleteSession) Filter() (res bson.M) {
	res = make(bson.M)
	res[fAccountId] = x.AccountId
	res[fAccountSessionData] = x.SessionId
	return
}

func (x IDeleteSession) Update() bson.M {
	return bson.M{
		mvars.OPull: bson.M{
			fAccountSessionData: x.SessionId,
		},
	}
}

func DeleteSessionOne(ctx context.Context, sessionId primitive.ObjectID) (err error) {
	var db, logger = vars.GetUtils(ctx)

	var res *mongo.UpdateResult
	if res, err = db.
		Collection(colAccount).
		UpdateOne(
			ctx,
			bson.M{
				fAccountSessionData: sessionId,
			},
			bson.M{
				mvars.OPull: bson.M{
					fAccountSessionData: sessionId,
				},
			},
		); err != nil {
		return
	}

	var accountId, isOk = res.UpsertedID.(primitive.ObjectID)
	if isOk {
		logger.Trace(
			"deleted session: accountId=%s, sessionId=%s",
			accountId.Hex(),
			sessionId.Hex(),
		)
	}

	return
}

func DeleteSessionAll(ctx context.Context, accountId primitive.ObjectID) (err error) {
	var db, logger = vars.GetUtils(ctx)
	if _, err = db.
		Collection(colAccount).
		UpdateOne(
			ctx,
			bson.M{
				fAccountId: accountId,
			},
			bson.M{
				mvars.OSet: bson.M{
					fAccountSessionData: make([]string, 0),
				},
			},
		); err != nil {
		return
	}

	logger.Trace("delete session all: accountId=%s", accountId.Hex())

	return
}
