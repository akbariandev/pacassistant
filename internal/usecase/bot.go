package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/akbariandev/pacassistant/config"
	"github.com/akbariandev/pacassistant/internal/domain"
	"github.com/akbariandev/pacassistant/pkg/cache"
	"github.com/akbariandev/pacassistant/pkg/logger"
	"github.com/akbariandev/pacassistant/transport/client"
	telegram "gopkg.in/telebot.v3"
	"math"
	"strconv"
)

const (
	StartMessage = "/start"
)

var (
	mainMenu          = &telegram.ReplyMarkup{}
	btnPacPrice       = mainMenu.Data("PAC Price", "pac_price")
	btnAddressBalance = mainMenu.Data("Address Balance", "address_balance")
)

type MessageStage int64

const (
	StageStart                MessageStage = 0
	StageGetAddressForBalance MessageStage = 1
)

type BotContext struct {
	Stage MessageStage
}

type BotUseCase struct {
	accountUsecase domain.AccountUsecase

	pactusClient *client.Pactus
	bot          *telegram.Bot

	priceCache cache.Cache[string, domain.Price]

	logger logger.Logger
}

func NewBot(
	accountUseCase domain.AccountUsecase,
	pactusClient *client.Pactus,
	bot *telegram.Bot,
	priceCache cache.Cache[string, domain.Price],
	logger logger.Logger,

) domain.BotUseCase {
	return &BotUseCase{
		accountUsecase: accountUseCase,
		logger:         logger,
		pactusClient:   pactusClient,
		bot:            bot,
		priceCache:     priceCache,
	}
}

func (b *BotUseCase) HandleMessages() {
	contextMap := make(map[int64]*BotContext)

	mainMenu.Inline(
		mainMenu.Row(btnPacPrice),
		mainMenu.Row(btnAddressBalance),
	)

	b.bot.Handle(StartMessage, func(c telegram.Context) error {
		return c.Send("Select an option", mainMenu)
	})

	b.bot.Handle(&btnPacPrice, func(c telegram.Context) error {
		priceData, ok := b.priceCache.Get(config.PriceCacheKey)
		if !ok {
			b.logger.Error(true, "failed to get price from cache")
			return nil
		}

		lastPrice, _ := strconv.ParseFloat(priceData.XeggexPacToUSDT.LastPrice, 64)
		yesterdayPrice, _ := strconv.ParseFloat(priceData.XeggexPacToUSDT.YesterdayPrice, 64)
		changePercentage := ((lastPrice - yesterdayPrice) / yesterdayPrice) * 100
		return c.Send(fmt.Sprintf("Current PAC Price: $%f    (%.f%%)", lastPrice, changePercentage))
	})

	b.bot.Handle(&btnAddressBalance, func(c telegram.Context) error {
		context := &BotContext{Stage: StageGetAddressForBalance}
		contextMap[c.Chat().ID] = context
		return c.Send("send your address")
	})

	b.bot.Handle(telegram.OnText, func(c telegram.Context) error {
		chatContext := contextMap[c.Message().Chat.ID]
		if chatContext == nil {
			return errors.New("context is nil")
		}

		switch chatContext.Stage {
		case StageGetAddressForBalance:
			account, err := b.getAccount(c.Message().Text)
			if err != nil {
				b.logger.Error(false, err.Error())
				chatContext.Stage = StageStart
				return c.Send("some errors happened in getting account balance. please try again later")
			}

			chatContext.Stage++
			balance := fmt.Sprintf("%.6f PAC", float64(account.Balance)*math.Pow10(-9))
			chatContext.Stage = StageStart
			return c.Send(fmt.Sprintf("Account Balance: %s", balance))
		default:
			return c.Send("Select an option", mainMenu)
		}
	})
}

func (b *BotUseCase) getAccount(address string) (domain.Account, error) {
	return b.accountUsecase.GetAccount(context.Background(), address)
}
