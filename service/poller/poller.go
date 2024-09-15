package poller

import (
	"context"
	"fmt"
	"github.com/heloayer/check-order-status/core/external_api/providers"
	"github.com/heloayer/check-order-status/internal/models"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type Poller struct {
	providerList []providers.Provider
	pollInterval time.Duration
	stopChan     chan struct{}
}

func NewPoller(providerList []providers.Provider, pollInterval time.Duration) *Poller {
	return &Poller{
		providerList: providerList,
		pollInterval: pollInterval,
		stopChan:     make(chan struct{}),
	}
}

func (p *Poller) Start() {
	log.Info().Msg("starting the poller")

	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), p.pollInterval)
			log.Info().Dur("poll_interval", p.pollInterval).Msg("polling providers")
			p.pollProviders(ctx)
			cancel()
		case <-p.stopChan:
			log.Info().Msg("stopping the poller")
			return
		}
	}
}

func (p *Poller) Stop() {
	log.Info().Msg("stopping the poller")
	close(p.stopChan)
}

func (p *Poller) pollProviders(ctx context.Context) {
	var wg sync.WaitGroup
	orderChan := make(chan models.Order)
	errChan := make(chan error)

	for _, provider := range p.providerList {
		wg.Add(1)
		go func(provider providers.Provider) {
			defer wg.Done()

			log.Info().Str("provider", fmt.Sprintf("%T", provider)).Msg("getting orders from provider")
			orders, err := provider.GetOrders(ctx)
			if err != nil {
				log.Error().Err(err).Msg("error getting orders from provider")
				errChan <- err
				return
			}
			for _, order := range orders {
				orderChan <- order
			}
		}(provider)
	}

	go func() {
		wg.Wait()
		close(orderChan)
		close(errChan)
	}()

	for {
		select {
		case order, ok := <-orderChan:
			if ok {
				log.Info().Msgf("Order ID: %s, Status: %s", order.ID, order.Status)
			} else {
				return
			}
		case err := <-errChan:
			log.Error().Err(err).Msg("error occurred while polling providers")
		case <-ctx.Done():
			log.Warn().Msg("polling canceled or timed out")
			return
		}
	}
}
