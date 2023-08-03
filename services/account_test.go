package services

import (
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/d3v-friends/pure-go/fnPanic"
	"github.com/d3v-friends/pure-go/fnReflect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestAccount(test *testing.T) {
	var tester = newTestTools()
	var accountSv = &AccountImpl{}
	var account *Account
	var err error

	fnSameAccount := func(t *testing.T, a *Account, b *Account) {
		assert.Equal(t, a.Id, b.Id)
		assert.Equal(t, true, reflect.DeepEqual(a.Identifier, b.Identifier))
		assert.Equal(t, true, reflect.DeepEqual(a.Permission, b.Permission))
		assert.Equal(t, true, reflect.DeepEqual(a.Property, b.Property))
		assert.Equal(t, true, reflect.DeepEqual(a.Verifier, b.Verifier))
		assert.Equal(t, a.CreatedAt, b.CreatedAt)
	}

	test.Run("create account", func(t *testing.T) {
		var ctx = tester.context()

		account, err = accountSv.Create(ctx, &ICreateAccount{
			Identifier: map[string]string{
				"email": gofakeit.Email(),
			},
			Property: map[string]string{
				"address": gofakeit.Address().Address,
			},
			Verifier: map[string]*IVerifier{
				"passwd": {
					Salt:   uuid.NewString(),
					Passwd: uuid.NewString(),
					Etc:    "passwd",
					Mode:   VerifyMode_COMPARE,
				},
			},
			Permission: map[string]bool{
				"signIn": true,
			},
		})

		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("account: %s", fnPanic.OnValue(json.Marshal(account)))
	})

	test.Run("readOne(Id)", func(t *testing.T) {
		var ctx = tester.context()
		var readAccount *Account
		readAccount, err = accountSv.ReadOne(ctx, &IReadAccount{
			Id: &account.Id,
		})

		if err != nil {
			t.Fatal(err)
		}
		fnSameAccount(t, account, readAccount)
	})

	test.Run("readOne(identifier)", func(t *testing.T) {
		var ctx = tester.context()
		var readAccount *Account
		readAccount, err = accountSv.ReadOne(ctx, &IReadAccount{
			Identifier: account.Identifier,
		})

		if err != nil {
			t.Fatal(err)
		}

		fnSameAccount(t, account, readAccount)
	})

	test.Run("readOne(property)", func(t *testing.T) {
		var ctx = tester.context()
		var readAccount *Account
		readAccount, err = accountSv.ReadOne(ctx, &IReadAccount{
			Property: account.Property,
		})

		if err != nil {
			t.Fatal(err)
		}

		fnSameAccount(t, account, readAccount)
	})

	test.Run("upsert", func(t *testing.T) {
		var ctx = tester.context()

		account.Identifier["email"] = gofakeit.Email()
		account.Permission["signIn"] = false
		account.Property["address"] = gofakeit.Address().Address

		var loadAccount = fnPanic.OnPointer(accountSv.Upsert(ctx, &IUpsertAccount{
			Filter: &IReadAccount{
				Id: fnReflect.ToPointer(account.Id),
			},
			Identifier: map[string]string{
				"email": account.Identifier["email"],
			},
			Property: map[string]string{
				"address": account.Property["address"],
			},
			Permission: map[string]bool{
				"signIn": account.Permission["signIn"],
			},
			Verifier: nil,
		}))

		assert.Equal(t, true, reflect.DeepEqual(account, loadAccount))

	})
}
