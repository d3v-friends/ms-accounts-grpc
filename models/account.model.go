package models

import (
	"fmt"
	"github.com/d3v-friends/mango/fn/fnBson"
	"github.com/d3v-friends/mango/mtype"
	"github.com/d3v-friends/mango/mvars"
	"github.com/d3v-friends/pure-go/fnParams"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// AccountModel Account 정의n
type (
	Account struct {
		Id        primitive.ObjectID `bson:"_id"`
		InTrx     bool               `bson:"inTrx"`
		Account   AccountAll         `bson:"account"`
		Session   []AccountSession   `bson:"session"`
		UpdatedAt time.Time          `bson:"updatedAt"`
		CreatedAt time.Time          `bson:"createdAt"`
		DeletedAt *time.Time         `bson:"deletedAt"`
	}

	AccountSession struct {
		Id        primitive.ObjectID `bson:"id"`
		CreatedAt time.Time          `bson:"createdAt"`
	}

	AccountAll struct {
		Data      AccountData    `bson:"data"`
		Histories []*AccountData `bson:"histories"`
	}

	AccountData struct {
		Identifier map[string]string   `bson:"identifier"`
		Property   map[string]string   `bson:"property"`
		Permission map[string]bool     `bson:"permission"`
		Verifier   map[string]Verifier `bson:"verifier"`
		CreatedAt  time.Time           `bson:"createdAt"`
	}

	Verifier struct {
		Salt   string `bson:"salt"`
		Passwd string `bson:"passwd"`
		Etc    string `bson:"etc"`
		Mode   string `bson:"mode"`
	}

	VerifierMap map[string]Verifier
)

const (
	colAccount = "accounts"

	fAccountId              = "_id"
	fAccountInTrx           = "inTrx"
	fAccountAccount         = "account"
	fAccountSession         = "session"
	fAccountSessionId       = "session.id"
	fAccountUpdatedAt       = "updatedAt"
	fAccountCreatedAt       = "createdAt"
	fAccountDeletedAt       = "deletedAt"
	fAccount2Data           = "account.data"
	fAccount2DataHistories  = "account.histories"
	fAccount2DataIdentifier = "account.data.identifier"
	fAccount2DataProperty   = "account.data.property"
	fAccount2DataPermission = "account.data.permission"
	fAccount2DataVerifier   = "account.data.verifier"
	fAccount2DataCreatedAt  = "account.data.createdAt"
)

var mgAccountModel = mtype.FnMigrateList{}

func NewAccountModel(data *AccountData) (res *Account) {
	now := time.Now()
	data.CreatedAt = now

	return &Account{
		Id:    primitive.NewObjectID(),
		InTrx: false,
		Account: AccountAll{
			Data:      *data,
			Histories: make([]*AccountData, 0),
		},
		Session:   make([]AccountSession, 0),
		UpdatedAt: now,
		CreatedAt: now,
	}
}

func (x Account) GetCollectionNm() string {
	return colAccount
}

func (x Account) GetMigrateList() mtype.FnMigrateList {
	return mgAccountModel
}

type IReadAccount struct {
	Id         *primitive.ObjectID
	Identifier map[string]string
	Property   map[string]string
	Permission map[string]bool
	SessionId  *primitive.ObjectID
	Sort       *ISortAccount
}

func (x *IReadAccount) Filter(IIsDelete ...bool) (res bson.D, err error) {
	res = make(bson.D, 0)

	if fnParams.Get(IIsDelete) {
		res = append(res, bson.E{
			Key: fAccountDeletedAt,
			Value: bson.M{
				"$ne": nil,
			},
		})
	}

	if x.Id != nil {
		res = append(res, bson.E{
			Key:   fAccountId,
			Value: *x.Id,
		})
	}

	if x.SessionId != nil {
		res = append(res, bson.E{
			Key:   fAccountSessionId,
			Value: *x.SessionId,
		})
	}

	if res, err = fnBson.MergeD(
		res,
		fnBson.ChangeMapToD(x.Identifier, fAccount2DataIdentifier),
	); err != nil {
		return
	}

	if res, err = fnBson.MergeD(
		res,
		fnBson.ChangeMapToD(x.Property, fAccount2DataProperty),
	); err != nil {
		return
	}

	if res, err = fnBson.MergeD(
		res,
		fnBson.ChangeMapToD(x.Permission, fAccount2DataPermission),
	); err != nil {
		return
	}

	return res, nil
}

func (x *IReadAccount) createSort() (bson.D, error) {
	if x.Sort == nil {
		return make(bson.D, 0), nil
	}
	return x.Sort.filter()
}

type ISortAccount struct {
	Identifier map[string]int64
	Property   map[string]int64
}

func (x *ISortAccount) filter() (res bson.D, err error) {
	res = make(bson.D, 0)

	if res, err = fnBson.MergeD(
		res,
		fnBson.ChangeMapToD(x.Identifier, fAccount2DataIdentifier),
	); err != nil {
		return
	}

	if res, err = fnBson.MergeD(
		res,
		fnBson.ChangeMapToD(x.Property, fAccount2DataProperty),
	); err != nil {
		return
	}

	return
}

type IDeleteAccountElem struct {
	Identifier []string
	Property   []string
	Verifier   []string
	Permission []string
}

func (x *IDeleteAccountElem) Change() (res bson.M, err error) {
	unset := make(bson.D, 0)

	if unset, err = fnBson.MergeD(
		unset,
		fnBson.ChangeEmptyToD(x.Identifier, fAccount2DataIdentifier),
	); err != nil {
		return
	}

	if unset, err = fnBson.MergeD(
		unset,
		fnBson.ChangeEmptyToD(x.Property, fAccount2DataProperty),
	); err != nil {
		return
	}

	if unset, err = fnBson.MergeD(
		unset,
		fnBson.ChangeEmptyToD(x.Verifier, fAccount2DataVerifier),
	); err != nil {
		return
	}

	if unset, err = fnBson.MergeD(
		unset,
		fnBson.ChangeEmptyToD(x.Permission, fAccount2DataPermission),
	); err != nil {
		return
	}

	return bson.M{
		mvars.OUnset: unset,
	}, nil
}

type IUpsertAccount struct {
	Identifier map[string]string
	Property   map[string]string
	Permission map[string]bool
	Verifier   VerifierMap
}

func (x *IUpsertAccount) Change() (_ bson.M, err error) {
	set := make(bson.D, 0)

	set = append(set, bson.E{
		Key:   fAccount2DataCreatedAt,
		Value: time.Now(),
	})

	if set, err = fnBson.MergeD(
		set,
		fnBson.ChangeMapToD(x.Identifier, fAccount2DataIdentifier),
	); err != nil {
		return
	}

	if set, err = fnBson.MergeD(
		set,
		fnBson.ChangeMapToD(x.Property, fAccount2DataProperty),
	); err != nil {
		return
	}

	if set, err = fnBson.MergeD(
		set,
		fnBson.ChangeMapToD(x.Permission, fAccount2DataPermission),
	); err != nil {
		return
	}

	if set, err = fnBson.MergeD(
		set,
		fnBson.ChangeMapToD(x.Identifier, fAccount2DataVerifier),
	); err != nil {
		return
	}

	return bson.M{
		mvars.OSet: set,
	}, nil
}

type IVerifyAccountSession struct {
	SessionId  primitive.ObjectID
	Permission map[string]bool
}

func (x IVerifyAccountSession) Filter() (res bson.M, err error) {
	res = make(bson.M)

	res[fAccountSessionId] = x.SessionId

	for key, value := range x.Permission {
		res[fmt.Sprintf("%s.%s", fAccount2DataPermission, key)] = value
	}

	return
}
