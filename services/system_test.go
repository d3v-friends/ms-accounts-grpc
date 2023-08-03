package services

import (
	"fmt"
	"github.com/d3v-friends/pure-go/fnMatcher"
	"github.com/d3v-friends/pure-go/fnReflect"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"testing"
	"time"
)

func TestSystem(test *testing.T) {
	tools := newTestTools()
	service := &SystemImpl{}

	test.Run("update", func(t *testing.T) {
		var err error
		ctx := tools.context()
		_, err = service.UpdateKeys(ctx, &IUpdateKeys{
			Identifier: []string{
				"email",
			},
			Property: []string{
				"address",
			},
			Permission: []string{
				"signIn",
			},
		})

		// 인덱싱에 약간의 시간이 필요
		<-time.NewTicker(time.Second * 3).C

		if err != nil {
			t.Fatal(err)
		}

		type Idx struct {
			Name   string `bson:"name"`
			Unique *bool  `bson:"unique"`
		}

		list := make([]Idx, 0)

		var cur *mongo.Cursor
		if cur, err = tools.DB.Collection("accounts").Indexes().List(ctx); err != nil {
			t.Fatal(err)
		}

		if err = cur.All(ctx, &list); err != nil {
			t.Fatal(err)
		}

		keyValueList := []struct {
			Name   string
			Unique *bool
		}{
			{
				Name: "_id_",
			},
			{
				Name:   "account.data.identifier.email_1",
				Unique: fnReflect.ToPointer(true),
			},
			{
				Name: "account.data.property.address_1",
			},
			{
				Name: "account.data.permission.signIn_1",
			},
			{
				Name:   "session.id_1",
				Unique: fnReflect.ToPointer(true),
			},
		}

		for _, v1 := range keyValueList {
			if !fnMatcher.Has(list, func(v Idx) bool {
				if v1.Name != v.Name {
					return false
				}

				if !reflect.DeepEqual(v1.Unique, v.Unique) {
					return false
				}

				return true
			}) {
				t.Fatalf("noHasIndex: %s", v1.Name)
			}
		}
	})

	test.Run("read", func(t *testing.T) {
		var ctx = tools.context()
		var index, err = service.ReadAccountIndex(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}

		if !fnMatcher.Has(index.Property, func(v string) bool {
			return v == "address"
		}) {
			err = fmt.Errorf("invalid property index")
		}

		if !fnMatcher.Has(index.Identifier, func(v string) bool {
			return v == "email"
		}) {
			err = fmt.Errorf("invalid identifer index")
		}

		if !fnMatcher.Has(index.Permission, func(v string) bool {
			return v == "signIn"
		}) {
			err = fmt.Errorf("invalid permission index")
		}

	})
}
