package usecase

import (
	"fmt"
	"github.com/akbariandev/pacassistant/config"
	"github.com/akbariandev/pacassistant/internal/domain"
	"github.com/akbariandev/pacassistant/pkg/cache"
	"github.com/akbariandev/pacassistant/pkg/logger"
	"github.com/akbariandev/pacassistant/transport/client"
	telegram "gopkg.in/telebot.v3"
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

type BotUseCase struct {
	pactusClient *client.Pactus
	bot          *telegram.Bot
	priceCache   cache.Cache[string, domain.Price]
	logger       logger.Logger
}

func NewBot(
	logger logger.Logger,
	pactusClient *client.Pactus,
	bot *telegram.Bot,
	priceCache cache.Cache[string, domain.Price],
) domain.BotUseCase {
	return &BotUseCase{
		logger:       logger,
		pactusClient: pactusClient,
		bot:          bot,
		priceCache:   priceCache,
	}
}

func (b *BotUseCase) HandleMessages() {

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

	/*	b.bot.Handle(&btnAddressBalance, func(c telegram.Context) error {
		c.Send("")
	})*/
}
