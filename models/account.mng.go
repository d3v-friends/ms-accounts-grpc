package models

import (
	"context"
	"github.com/d3v-friends/mango/mvars"
	"github.com/d3v-friends/ms-accounts-grpc/vars"
	"github.com/d3v-friends/pure-go/fnReflect"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func CreateAccount(ctx context.Context, i *AccountData) (account *Account, err error) {
	db, logger := vars.GetUtils(ctx)

	account = NewAccountModel(i)

	if _, err = db.
		Collection(account.GetCollectionNm()).
		InsertOne(ctx, account); err != nil {
		return
	}

	logger.Trace("created user: id=%s", account.Id.Hex())

	return account, nil
}

func ReadAccountOne(ctx context.Context, i *IReadAccount) (account *Account, err error) {
	db, logger := vars.GetUtils(ctx)

	var filter bson.D
	if filter, err = i.Filter(); err != nil {
		return
	}

	var sort bson.D
	if sort, err = i.createSort(); err != nil {
		return
	}

	var opt *options.FindOneOptions
	if 0 < len(sort) {
		opt = &options.FindOneOptions{
			Sort: sort,
		}
	}

	var singleRes *mongo.SingleResult
	if singleRes = db.Collection(colAccount).FindOne(ctx, filter, opt); singleRes.Err() != nil {
		err = singleRes.Err()
		return
	}

	account = new(Account)
	if err = singleRes.Decode(account); err != nil {
		return
	}

	logger.Trace("found user: id=%s", account.Id.Hex())

	return
}

func ReadAccountAll(ctx context.Context, i *IReadAccount) (ls []Account, err error) {

	var db, logger = vars.GetUtils(ctx)
	var filter bson.D
	ls = make([]Account, 0)

	if filter, err = i.Filter(); err != nil {
		return
	}

	var sort bson.D
	if sort, err = i.createSort(); err != nil {
		return
	}

	var opt *options.FindOptions
	if 0 < len(sort) {
		opt = &options.FindOptions{
			Sort: sort,
		}
	}

	var cur *mongo.Cursor
	if cur, err = db.Collection(colAccount).Find(ctx, filter, opt); err != nil {
		return
	}

	if err = cur.All(ctx, &ls); err != nil {
		return
	}

	logger.Trace("found users: len(ls)=%d", len(ls))

	return
}

func ReadAccountList(ctx context.Context, i *IReadAccount, p *IPager) (list *List[Account], err error) {
	var db, logger = vars.GetUtils(ctx)
	var filter bson.D
	if filter, err = i.Filter(); err != nil {
		return
	}

	var sort bson.D
	if sort, err = i.createSort(); err != nil {
		return
	}

	list = &List[Account]{
		Page:  p.Page,
		Size:  p.Size,
		Total: 0,
		List:  make([]Account, 0),
	}

	if list.Total, err = db.Collection(colAccount).CountDocuments(ctx, filter); err != nil {
		return
	}

	var opt = &options.FindOptions{
		Skip:  p.Skip(),
		Limit: p.Limit(),
	}
	if 0 < len(sort) {
		opt.Sort = sort
	}

	var cur *mongo.Cursor
	if cur, err = db.Collection(colAccount).Find(ctx, filter, opt); err != nil {
		return
	}

	if err = cur.All(ctx, &list.List); err != nil {
		return
	}

	logger.Trace("found user list: len(ls)=%d", len(list.List))

	return
}

func UpsertAccount(
	ctx context.Context,
	i *IReadAccount,
	d *IUpsertAccount,
) (account *Account, err error) {
	db, logger := vars.GetUtils(ctx)

	var update bson.M
	if update, err = d.Change(); err != nil {
		return
	}

	var filter bson.D
	if filter, err = i.Filter(); err != nil {
		return
	}

	var sort bson.D
	if sort, err = i.createSort(); err != nil {
		return
	}

	var opt = &options.FindOneAndUpdateOptions{
		Upsert: fnReflect.ToPointer(true),
	}

	if 0 < len(sort) {
		opt.Sort = sort
	}

	var cur *mongo.SingleResult
	if cur = db.Collection(colAccount).FindOneAndUpdate(
		ctx,
		filter,
		update,
		opt,
	); cur.Err() != nil {
		err = cur.Err()
		return
	}

	account = new(Account)
	if err = cur.Decode(account); err != nil {
		return
	}

	if _, err = db.Collection(colAccount).UpdateOne(ctx, filter, &bson.M{
		mvars.OPush: &bson.M{
			fAccount2DataHistories: account.Account.Data,
		},
	}); err != nil {
		return
	}

	logger.Trace(
		"upsert user: id=%s, upsertData=%d",
		account.Id.Hex(),
		len(filter),
	)

	return ReadAccountOne(ctx, i)
}

func DeleteAccount(ctx context.Context, i *IReadAccount) (err error) {
	var db, logger = vars.GetUtils(ctx)

	var filter bson.D
	if filter, err = i.Filter(); err != nil {
		return
	}

	var res *mongo.UpdateResult
	if res, err = db.
		Collection(colAccount).
		UpdateOne(ctx, filter, bson.M{
			mvars.OSet: bson.M{
				fAccountDeletedAt: time.Now(),
			},
		}); err != nil {
		return
	}

	var accountId, isOk = res.UpsertedID.(primitive.ObjectID)
	if isOk {
		logger.Trace("deleted account: id=%s", accountId.Hex())
	}

	return
}

func DeleteAccountElem(
	ctx context.Context,
	i *IReadAccount,
	elem *IDeleteAccountElem,
) (account *Account, err error) {
	var db, logger = vars.GetUtils(ctx)
	var filter bson.D
	if filter, err = i.Filter(); err != nil {
		return
	}

	var update bson.M
	if update, err = elem.Change(); err != nil {
		return
	}

	var res *mongo.UpdateResult
	if res, err = db.
		Collection(colAccount).
		UpdateOne(ctx, filter, update); err != nil {
	}

	var accountId, isOk = res.UpsertedID.(primitive.ObjectID)
	if isOk {
		logger.Trace("deleted account elem: id=%s", accountId.Hex())
	}

	return ReadAccountOne(ctx, i)
}
