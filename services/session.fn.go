package services

import (
	"github.com/d3v-friends/ms-accounts-grpc/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (x *ICreateSession) CreateSession() *models.ICreateSession {
	return &models.ICreateSession{
		Identifier: x.Identifier,
		Permission: x.Permission,
		Verifier:   x.Verifier,
	}
}

func newSession(i *models.AccountSession) (*Session, error) {
	return &Session{
		SessionId: i.Id.Hex(),
		CreatedAt: i.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (x *IVerifySession) IVerifyAccountSession() (res *models.IVerifyAccountSession, err error) {
	var sessionId primitive.ObjectID
	if sessionId, err = primitive.ObjectIDFromHex(x.SessionId); err != nil {
		return
	}

	return &models.IVerifyAccountSession{
		SessionId:  sessionId,
		Permission: x.Permission,
	}, nil
}

func (x *IDeleteSessionOne) ParseSessionId() (res primitive.ObjectID, err error) {
	return primitive.ObjectIDFromHex(x.SessionId)
}

func (x *IDeleteSessionAll) ParseAccountId() (res primitive.ObjectID, err error) {
	return primitive.ObjectIDFromHex(x.AccountId)
}
