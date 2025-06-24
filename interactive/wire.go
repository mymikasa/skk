//go:build wireinject

package main

import (
	"github.com/mymikasa/skk/interactive/events"
	"github.com/mymikasa/skk/interactive/grpc"
	"github.com/mymikasa/skk/interactive/ioc"
	repository2 "github.com/mymikasa/skk/interactive/repository"
	cache2 "github.com/mymikasa/skk/interactive/repository/cache"
	dao2 "github.com/mymikasa/skk/interactive/repository/dao"
	service2 "github.com/mymikasa/skk/interactive/service"

	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(ioc.InitDB,
	ioc.InitLogger,
	ioc.InitSaramaClient,
	ioc.InitRedis)

var interactiveSvcSet = wire.NewSet(dao2.NewGORMInteractiveDAO,
	cache2.NewInteractiveRedisCache,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService,
)

func InitApp() *App {
	wire.Build(thirdPartySet,
		interactiveSvcSet,
		grpc.NewInteractiveServiceServer,
		events.NewInteractiveReadEventConsumer,
		ioc.InitConsumers,
		ioc.NewGrpcxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
