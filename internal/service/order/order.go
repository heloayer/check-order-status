package order

import (
	"context"
	"fmt"
	"github.com/heloayer/check-order-status/core/external_api/providers"
	"github.com/heloayer/check-order-status/internal/models"
	"github.com/heloayer/check-order-status/internal/resources/http/handler/dto"
	"github.com/rs/zerolog/log"
	"sync"
)

type OrderService struct {
	providerList []providers.Provider
}

func NewOrderService(providerList []providers.Provider) OrderService {
	return OrderService{
		providerList: providerList,
	}
}

type Order interface {
	GetOrderByID(ctx context.Context, id string) (models.Order, error)
	CreateOrder(ctx context.Context, request dto.OrderRequest, provider string) error
}

func (ord *OrderService) GetOrderByID(ctx context.Context, id string) (models.Order, error) {
	log.Info().Str("order_id", id).Msg("starting to find order by ID")

	var wg sync.WaitGroup
	orderChan := make(chan models.Order)
	errChan := make(chan error)
	for _, prov := range ord.providerList {
		wg.Add(1)
		go func(provider providers.Provider) {
			defer wg.Done()
			orders, err := provider.GetOrders(ctx)
			if err != nil {
				log.Error().Err(err).Msg("error getting orders from provider")
				errChan <- err
				return
			}
			for _, order := range orders {
				if order.ID == id {
					log.Info().Str("order_id", id).Msg("order found")
					orderChan <- order
					return
				}
			}
		}(prov)
	}
	go func() {
		wg.Wait()
		close(orderChan)
		close(errChan)
	}()

	select {
	case order := <-orderChan:
		log.Info().Str("order_id", id).Msg("successfully found order")
		return order, nil
	case err := <-errChan:
		log.Error().Err(err).Msg("failed to find order")
		return models.Order{}, fmt.Errorf("couldn't find order: %w", err)
	case <-ctx.Done():
		log.Warn().Msg("context canceled or timed out")
		return models.Order{}, ctx.Err()
	}
}

func (o *OrderService) CreateOrder(ctx context.Context, request dto.OrderRequest, provider string) error {
	log.Info().Str("provider", provider).Msg("starting to create order")

	var selectedProvider providers.Provider

	for _, prov := range o.providerList {
		switch provider {
		case models.Glovo:
			if jsonProv, ok := prov.(*providers.JsonProvider); ok {
				selectedProvider = jsonProv
				break
			}
		case models.Wolt:
			if mockProv, ok := prov.(*providers.MockProvider); ok {
				selectedProvider = mockProv
				break
			}
		}
	}
	if selectedProvider == nil {
		log.Error().Str("provider", provider).Msg("invalid provider specified")
		return fmt.Errorf("invalid provider specified")
	}

	return selectedProvider.CreateOrder(ctx, request.ToModel())
}
