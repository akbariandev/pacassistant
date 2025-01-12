package job

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/akbariandev/pacassistant/config"
	"github.com/akbariandev/pacassistant/internal/domain"
	"github.com/akbariandev/pacassistant/pkg/cache"
	"github.com/akbariandev/pacassistant/pkg/logger"
	"net/http"
	"sync"
	"time"
)

const (
	_defaultXeggexPriceEndpoint   = "https://api.xeggex.com/api/v2/market/getbysymbol/Pactus%2Fusdt"
	_defaultExbitronPriceEndpoint = "https://api.exbitron.digital/api/v1/cg/tickers"
)

type price struct {
	cache  cache.Cache[string, domain.Price]
	ticker *time.Ticker
	logger logger.Logger
	ctx    context.Context
	cancel context.CancelFunc
}

func NewPrice(
	cache cache.Cache[string, domain.Price],
	logger logger.Logger,
) Job {
	ctx, cancel := context.WithCancel(context.Background())
	return &price{
		cache:  cache,
		ticker: time.NewTicker(128 * time.Second),
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (p *price) Start() {
	p.start()
	go p.runTicker()
}

func (p *price) start() {
	var (
		wg       sync.WaitGroup
		price    domain.Price
		xeggex   domain.XeggexPriceResponse
		exbitron domain.ExbitronPriceResponse
	)

	ctx := context.Background()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.getPrice(ctx, _defaultXeggexPriceEndpoint, &xeggex); err != nil {
			p.logger.Error(true, err.Error())
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.getPrice(ctx, _defaultExbitronPriceEndpoint, &exbitron); err != nil {
			p.logger.Error(true, err.Error())
		}
	}()

	wg.Wait()

	price.XeggexPacToUSDT = xeggex
	price.ExbitronPacToUSDT = exbitron

	ok := p.cache.Exists(config.PriceCacheKey)
	if ok {
		p.cache.Update(config.PriceCacheKey, price, 0)
	} else {
		p.cache.Add(config.PriceCacheKey, price, 0)
	}
}

func (p *price) runTicker() {
	for {
		select {
		case <-p.ctx.Done():
			return

		case <-p.ticker.C:
			p.start()
		}
	}
}

func (p *price) getPrice(ctx context.Context, endpoint string, vPtrResp any) error {
	cli := http.DefaultClient

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response code is %v", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	return dec.Decode(vPtrResp)
}

func (p *price) Stop() {
	p.ticker.Stop()
}
