package main

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"go-tour/config"
	"go-tour/internal/dao"
	"go-tour/internal/must"
	"go-tour/internal/services"
	"go-tour/migration"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	var ctx = context.TODO()
	cfg := config.ReadConfigAndArg()

	logger, sentry, err := must.NewLogger(cfg.SentryDSN, cfg.ServiceName+"-app")
	if err != nil {
		log.Fatalf("logger: %v", err)
	}

	defer logger.Sync()
	defer sentry.Flush(2 * time.Second)

	db := must.ConnectDb(cfg.Db)
	err = migration.Migration(db)
	if err != nil {
		log.Fatalf("migration: %v", err)
	}

	if err := migration.AutoSeedingData(db); err != nil {
		//log.Fatalf("seeding: %v", err)
	}

	_, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))

	//dao
	userDao := dao.NewUser(db)
	middlewareAuth := NewMiddleware(cfg.AuthenticationPubSecretKey)
	opt := []grpc.ServerOption{
		grpc.StreamInterceptor(auth.StreamServerInterceptor(middlewareAuth.AuthMiddleware)),
		grpc.UnaryInterceptor(auth.UnaryServerInterceptor(middlewareAuth.AuthMiddleware)),
	}

	must.NewServer(ctx, cfg,
		opt,
		services.NewUserService(
			logger,
			cfg,
			db,
			userDao,
		),
		services.NewAuthenticationService(
			logger,
			cfg,
			userDao,
		),
	)
}
