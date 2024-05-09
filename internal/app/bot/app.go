package bot

import (
	"context"
	"errors"
	telegram "gopkg.in/telebot.v3"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/akbariandev/pacassistant/internal/domain"
	jobPkg "github.com/akbariandev/pacassistant/internal/job"

	"github.com/akbariandev/pacassistant/pkg/cache"
	"github.com/akbariandev/pacassistant/pkg/scheduler"

	"github.com/akbariandev/pacassistant/transport/client"

	"github.com/akbariandev/pacassistant/pkg/logger"

	"github.com/akbariandev/pacassistant/internal/usecase"

	"github.com/akbariandev/pacassistant/config"
	"github.com/akbariandev/pacassistant/pkg/mongodb"
)

type Bootstrapper interface {
	Run(ctx context.Context)
	Migration(ctx context.Context) error
	Shutdown(ctx context.Context) error
	GetServiceConfig() *config.Config[config.ExtraData]
	GetMongodbConnector() mongodb.Connector
}

type Application struct {
	pactusClient     *client.Pactus
	telegramBot      *telegram.Bot
	serviceConfig    *config.Config[config.ExtraData]
	mongodbConnector mongodb.Connector
	logger           logger.Logger
	transaction      *mongodb.Transaction
	scheduler        scheduler.Scheduler
}

func New(
	ctx context.Context,
	pactusClient *client.Pactus,
	telegramBot *telegram.Bot,
	serviceConfig *config.Config[config.ExtraData],
	logger logger.Logger,
) (Bootstrapper, error) {
	app := new(Application)

	if serviceConfig == nil {
		return nil, errors.New("service config is nil")
	}

	/*	database, err := newMongodb(ctx, serviceConfig.Database.Mongodb.URI, serviceConfig.Database.Mongodb.DatabaseName)
		if err != nil {
			return nil, err
		}
	*/
	app.serviceConfig = serviceConfig
	app.logger = logger
	//app.mongodbConnector = database
	app.pactusClient = pactusClient
	app.telegramBot = telegramBot
	//app.transaction = mongodb.NewTransaction(database.GetMongoClient())
	app.scheduler = scheduler.NewScheduler()

	return app, nil
}

func (a *Application) Run(ctx context.Context) {
	a.registerHandlerLayer()
	go a.scheduler.Run()
	a.logger.InfoContext(ctx, false, "Scheduler Started Successfully!")
	go a.telegramBot.Start()
	a.logger.InfoContext(ctx, false, "Telegram Bot Started Successfully!")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {}
}

func (a *Application) Shutdown(ctx context.Context) error {
	return nil
}

func (a *Application) GetServiceConfig() *config.Config[config.ExtraData] {
	return a.serviceConfig
}

func (a *Application) GetMongodbConnector() mongodb.Connector {
	return a.mongodbConnector
}

func (a *Application) registerHandlerLayer() {
	priceCache := cache.NewBasic[string, domain.Price](0 * time.Second)

	accountUseCase := usecase.NewAccount(a.pactusClient, a.logger, nil)
	botUseCase := usecase.NewBot(accountUseCase, a.pactusClient, a.telegramBot, priceCache, a.logger)
	botUseCase.HandleMessages()

	// background jobs
	priceJob := jobPkg.NewPrice(priceCache, a.logger)
	a.scheduler.Submit(priceJob)
}

func (a *Application) Migration(ctx context.Context) error {
	return nil
}

func newMongodb(ctx context.Context, uri string, databaseName string) (mongodb.Connector, error) {
	connector, err := mongodb.New(ctx, uri, &sync.Mutex{})
	if err != nil {
		return nil, err
	}

	connector.SetDatabase(databaseName)

	return connector, nil
}
