package services

import (
	"github.com/d3v-friends/ms-accounts-grpc/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (x *ICreateAccount) changeVerifier() models.VerifierMap {
	return changeVerifierMap(x.Verifier)
}

func (x *IReadAccount) changeReadAccount() (res *models.IReadAccount, err error) {
	res = &models.IReadAccount{
		Identifier: x.Identifier,
		Property:   x.Property,
		Permission: x.Permission,
	}

	if x.Id != nil {
		var id primitive.ObjectID
		if id, err = primitive.ObjectIDFromHex(*x.Id); err != nil {
			return
		}
		res.Id = &id
	}

	if x.Sort != nil {
		res.Sort = &models.ISortAccount{
			Identifier: x.Sort.Identifier,
			Property:   x.Sort.Property,
		}
	}

	return
}

func (x *IUpsertAccount) changeUpsertAccount() *models.IUpsertAccount {
	return &models.IUpsertAccount{
		Identifier: x.Identifier,
		Property:   x.Property,
		Permission: x.Permission,
		Verifier:   changeVerifierMap(x.Verifier),
	}
}

func (x *IReadAccountList) pager() *models.IPager {
	return &models.IPager{
		Page: x.Page,
		Size: x.Size,
	}
}

func (x *IDeleteAccountElem) changeDeleteAccount() *models.IDeleteAccountElem {
	return &models.IDeleteAccountElem{
		Identifier: x.Identifier,
		Property:   x.Property,
		Verifier:   x.Verifier,
		Permission: x.Permission,
	}
}

func newAccount(v *models.Account) (res *Account) {
	res = &Account{
		Id:         v.Id.Hex(),
		Identifier: v.Account.Data.Identifier,
		Verifier:   newVerifier(v.Account.Data.Verifier),
		Property:   v.Account.Data.Property,
		Permission: v.Account.Data.Permission,
		CreatedAt:  v.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  v.Account.Data.CreatedAt.Format(time.RFC3339),
	}
	return
}

func newVerifier(v models.VerifierMap) (res map[string]*Verifier) {
	res = make(map[string]*Verifier)
	for key, verifier := range v {
		res[key] = &Verifier{
			Salt: verifier.Salt,
			Mode: verifier.Mode,
		}
	}
	return
}

func newAccountAll(v []models.Account) *AccountAll {
	return &AccountAll{
		List: changeAccountList(v),
	}
}

func newAccountList(v *models.List[models.Account]) *AccountList {
	return &AccountList{
		Page:  v.Page,
		Size:  v.Size,
		Total: v.Total,
		List:  changeAccountList(v.List),
	}
}

func changeAccountList(v []models.Account) (res []*Account) {
	for _, account := range v {
		res = append(res, newAccount(&account))
	}
	return
}

func changeVerifier(v *IVerifier) models.Verifier {
	return models.Verifier{
		Salt:   v.Salt,
		Passwd: v.Passwd,
		Etc:    v.Etc,
		Mode:   v.Mode.String(),
	}
}

func changeVerifierMap(v map[string]*IVerifier) (res models.VerifierMap) {
	res = make(models.VerifierMap)
	for key, verifier := range v {
		res[key] = changeVerifier(verifier)
	}
	return
}
