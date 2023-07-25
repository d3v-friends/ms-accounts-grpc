package main

import (
	"context"
	"github.com/d3v-friends/go-grpc/fn/fnGrpc"
	"github.com/d3v-friends/mango"
	"github.com/d3v-friends/ms-accounts-grpc/models"
	"github.com/d3v-friends/ms-accounts-grpc/services"
	"github.com/d3v-friends/ms-accounts-grpc/vars"
	"github.com/d3v-friends/pure-go/fnEnv"
	"github.com/d3v-friends/pure-go/fnLogger"
	"github.com/d3v-friends/pure-go/fnPanic"
	"github.com/google/uuid"
	cronV3 "github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func main() {
	var logger = fnLogger.NewDefaultLogger()
	var mClient = fnPanic.OnValue(mango.NewClient(&mango.ClientOpt{
		Host:     fnEnv.Read("MG_HOST"),
		Username: fnEnv.Read("MG_USERNAME"),
		Password: fnEnv.Read("MG_PASSWORD"),
		Database: fnEnv.Read("MG_DATABASE"),
	}))
	fnPanic.On(mClient.Migrate(context.TODO(), models.All...))

	var cron = cronV3.New()
	// todo cron 기능 이곳에 추가
	go cron.Run()
	var server = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(getInterceptor(mClient.Database(), logger)),
	)

	services.RegisterAccountsServer(server, &services.AccountImpl{})
	services.RegisterSystemsServer(server, &services.SystemImpl{})
	services.RegisterSessionsServer(server, &services.SessionImpl{})

	fnPanic.On(fnGrpc.Listen(server, fnEnv.Read("PORT")))
}

func getInterceptor(db *mongo.Database, logger fnLogger.IfLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		requestAt := time.Now()
		requestLogger := logger.WithFields(fnLogger.Fields{
			"requestId": uuid.NewString(),
			"method":    info.FullMethod,
			"requestAt": time.Now(),
		})
		requestLogger.Trace("requested")

		ctx = vars.SetUtils(ctx, db, logger)

		defer func() {
			responseAt := time.Now()
			responseLogger := requestLogger.
				WithFields(fnLogger.Fields{
					"responseAt": responseAt,
					"durations":  responseAt.UnixMilli() - requestAt.UnixMilli(),
				})

			if err == nil {
				responseLogger.
					Trace("responded")
			} else {
				responseLogger.
					Error("error: err=%s", err.Error())
			}
		}()

		return handler(ctx, req)
	}
}