package vars

import (
	"context"
	"fmt"
	"github.com/d3v-friends/pure-go/fnLogger"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	CtxMongoDB = "CTX_MONGO_DB"
	CtxLogger  = "CTX_LOGGER"
)

func GetDB(ctx context.Context) (res *mongo.Database) {
	var valid bool
	if res, valid = ctx.Value(CtxMongoDB).(*mongo.Database); !valid {
		err := fmt.Errorf("not found mongo database in context")
		panic(err)
	}
	return
}

func SetDB(ctx context.Context, db *mongo.Database) context.Context {
	return context.WithValue(ctx, CtxMongoDB, db)

}

func GetLogger(ctx context.Context) (res fnLogger.IfLogger) {
	var valid bool
	if res, valid = ctx.Value(CtxLogger).(fnLogger.IfLogger); !valid {
		err := fmt.Errorf("not found fnLogger.IfLogger")
		panic(err)
	}
	return
}

func SetLogger(ctx context.Context, logger fnLogger.IfLogger) context.Context {
	return context.WithValue(ctx, CtxLogger, logger)
}

func GetUtils(ctx context.Context) (db *mongo.Database, logger fnLogger.IfLogger) {
	return GetDB(ctx), GetLogger(ctx)
}

func SetUtils(ctx context.Context, db *mongo.Database, logger fnLogger.IfLogger) context.Context {
	ctx = SetDB(ctx, db)
	ctx = SetLogger(ctx, logger)
	return ctx
}
