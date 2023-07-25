package services

import (
	"context"
	"github.com/d3v-friends/ms-accounts-grpc/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	SessionImpl struct{}
)

func (x *SessionImpl) Create(ctx context.Context, i *ICreateSession) (res *Session, err error) {
	var session *models.SessionData
	if session, err = models.
		CreateSession(ctx, i.CreateSession()); err != nil {
		return
	}
	return newSession(session)
}

func (x *SessionImpl) Verify(ctx context.Context, i *IVerifySession) (res *Account, err error) {
	var iVerifySession *models.IVerifyAccountSession
	if iVerifySession, err = i.IVerifyAccountSession(); err != nil {
		return
	}

	var account *models.Account
	if account, err = models.VerifySession(ctx, iVerifySession); err != nil {
		return
	}

	res = newAccount(account)
	return
}

func (x *SessionImpl) DeleteOne(ctx context.Context, i *IDeleteSessionOne) (_ *Empty, err error) {
	var sessionId primitive.ObjectID
	if sessionId, err = i.ParseSessionId(); err != nil {
		return
	}

	if err = models.DeleteSessionOne(ctx, sessionId); err != nil {
		return
	}

	return
}

func (x *SessionImpl) DeleteAll(ctx context.Context, i *IDeleteSessionAll) (_ *Empty, err error) {
	var accountId primitive.ObjectID
	if accountId, err = i.ParseAccountId(); err != nil {
		return
	}

	if err = models.DeleteSessionAll(ctx, accountId); err != nil {
		return
	}

	return
}

func (x *SessionImpl) mustEmbedUnimplementedSessionsServer() {
}
