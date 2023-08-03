package services

import (
	"fmt"
	"github.com/d3v-friends/pure-go/fnPanic"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestSession(test *testing.T) {
	var tester = newTestTools()
	var sessionSv = &SessionImpl{}
	var accountSv = &AccountImpl{}
	var systemSv = &SystemImpl{}
	var account = fnPanic.OnPointer(accountSv.Create(tester.context(), tester.createAccount()))

	fnPanic.On(tester.indexAccount(systemSv))

	test.Run("create session", func(t *testing.T) {
		var ctx = tester.context()

		var session = fnPanic.OnPointer(sessionSv.Create(ctx, tester.iCreateSession(account)))

		fmt.Printf("sessionId = %s\n", session.SessionId)
	})

	test.Run("create session and verify", func(t *testing.T) {
		var ctx = tester.context()
		var session = fnPanic.OnPointer(sessionSv.Create(ctx, tester.iCreateSession(account)))

		var loadAccount = fnPanic.OnPointer(sessionSv.Verify(ctx, &IVerifySession{
			SessionId: session.SessionId,
			Permission: map[string]bool{
				"signIn": true,
			},
		}))

		assert.Equal(t, true, reflect.DeepEqual(account, loadAccount))
	})

	test.Run("verify and delete", func(t *testing.T) {
		var err error
		var ctx = tester.context()
		var session = fnPanic.OnPointer(sessionSv.Create(ctx, tester.iCreateSession(account)))

		fnPanic.OnPointer(sessionSv.DeleteOne(ctx, &IDeleteSessionOne{
			SessionId: session.SessionId,
		}))

		_, err = sessionSv.Verify(ctx, &IVerifySession{
			SessionId: session.SessionId,
			Permission: map[string]bool{
				"signIn": true,
			},
		})

		if err == nil {
			t.Fatalf("session is not deleted: sessionId=%s", session.SessionId)
		}
	})

}
