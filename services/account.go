package services

import (
	"context"
	"github.com/d3v-friends/ms-accounts-grpc/models"
	"time"
)

type AccountImpl struct {
}

func (x *AccountImpl) Create(ctx context.Context, i *ICreateAccount) (res *Account, err error) {
	var account *models.Account
	if account, err = models.CreateAccount(ctx, &models.AccountData{
		Identifier: i.Identifier,
		Property:   i.Property,
		Permission: i.Permission,
		Verifier:   i.changeVerifier(),
		CreatedAt:  time.Now(),
	}); err != nil {
		return
	}

	res = newAccount(account)

	return
}

func (x *AccountImpl) ReadOne(ctx context.Context, i *IReadAccount) (res *Account, err error) {
	var iReadAccount *models.IReadAccount
	if iReadAccount, err = i.changeReadAccount(); err != nil {
		return
	}

	var account *models.Account
	if account, err = models.ReadAccountOne(ctx, iReadAccount); err != nil {
		return nil, err
	}

	res = newAccount(account)
	return
}

func (x *AccountImpl) ReadAll(ctx context.Context, i *IReadAccount) (res *AccountAll, err error) {
	var iReadAccount *models.IReadAccount
	if iReadAccount, err = i.changeReadAccount(); err != nil {
		return
	}

	var accountAll []models.Account
	if accountAll, err = models.ReadAccountAll(ctx, iReadAccount); err != nil {
		return nil, err
	}

	res = newAccountAll(accountAll)
	return
}

func (x *AccountImpl) ReadList(ctx context.Context, i *IReadAccountList) (res *AccountList, err error) {
	var iReadAccount *models.IReadAccount
	if iReadAccount, err = i.Filter.changeReadAccount(); err != nil {
		return
	}

	var accountList *models.List[models.Account]
	if accountList, err = models.ReadAccountList(
		ctx,
		iReadAccount,
		i.pager(),
	); err != nil {
		return nil, err
	}

	res = newAccountList(accountList)

	return
}

func (x *AccountImpl) Upsert(ctx context.Context, i *IUpsertAccount) (res *Account, err error) {
	var iReadAccount *models.IReadAccount
	if iReadAccount, err = i.Filter.changeReadAccount(); err != nil {
		return
	}

	var account *models.Account
	if account, err = models.UpsertAccount(
		ctx,
		iReadAccount,
		i.changeUpsertAccount(),
	); err != nil {
		return
	}

	res = newAccount(account)
	return
}

func (x *AccountImpl) DeleteElem(ctx context.Context, i *IDeleteAccountElem) (res *Account, err error) {
	var iReadAccount *models.IReadAccount
	if iReadAccount, err = i.Filter.changeReadAccount(); err != nil {
		return
	}

	var iDeleteAccount = i.changeDeleteAccount()

	var account *models.Account
	if account, err = models.DeleteAccountElem(
		ctx,
		iReadAccount,
		iDeleteAccount,
	); err != nil {
		return
	}

	res = newAccount(account)
	return
}

func (x *AccountImpl) mustEmbedUnimplementedAccountsServer() {
}
