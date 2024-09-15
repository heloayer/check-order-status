package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/heloayer/check-order-status/config"
	"github.com/heloayer/check-order-status/internal/models"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

type Provider interface {
	GetOrders(ctx context.Context) ([]models.Order, error)
	CreateOrder(ctx context.Context, order models.Order) error
}

type JsonProvider struct {
	endpoint string
}

type MockProvider struct {
	endpoint string
}

func NewJSONProvider(cfg config.ProviderConfig) Provider {
	return &JsonProvider{
		endpoint: cfg.Endpoint,
	}
}

func NewMockAPIProvider(cfg config.ProviderConfig) Provider {
	return &MockProvider{
		endpoint: cfg.Endpoint,
	}
}

func (p *JsonProvider) GetOrders(ctx context.Context) ([]models.Order, error) {
	req, err := http.NewRequest("GET", p.endpoint+"/orders", nil)
	if err != nil {
		log.Error().Err(err).Msg("error creating http request")
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Error().Err(err).Msg("error making get request to JSON Server orders endpoint")
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("failed to read response body")
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		log.Error().Str("status", resp.Status).Str("body", string(body)).Msg("failed to get orders")
		return nil, fmt.Errorf("failed to get orders: %s", resp.Status)
	}

	var orders []models.Order
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		log.Error().Err(err).Msg("error decoding json response")
		return nil, err
	}

	log.Info().Int("count", len(orders)).Msg("successfully retrieved Glovo orders")
	return orders, nil
}

func (p *MockProvider) GetOrders(ctx context.Context) ([]models.Order, error) {
	req, err := http.NewRequest("GET", p.endpoint+"/v1/orders", nil)
	if err != nil {
		log.Error().Err(err).Msg("error creating http request")
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Error().Err(err).Msg("error making get request to Mock API orders endpoint")
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("failed to read response body")
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		log.Error().Str("status", resp.Status).Str("body", string(body)).Msg("failed to get orders")
		return nil, fmt.Errorf("failed to get orders: %s, body: %s", resp.Status, string(body))
	}

	var orders []models.Order
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		log.Error().Err(err).Msg("error decoding json response")
		return nil, err
	}

	log.Info().Int("count", len(orders)).Msg("successfully retrieved WOLT orders")
	return orders, nil
}

func (p *JsonProvider) CreateOrder(ctx context.Context, order models.Order) error {

	log.Info().Msg("starting to create Glovo order")

	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Error().Err(err).Msg("error marshaling order to json")
		return err
	}

	req, err := http.NewRequest("POST", p.endpoint+"/orders", bytes.NewBuffer(orderJSON))
	if err != nil {
		log.Error().Err(err).Msg("error creating http request")
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Error().Err(err).Msg("error making http request")
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("failed to read response body")
			return fmt.Errorf("failed to read response body: %w", err)
		}
		log.Error().
			Int("status", resp.StatusCode).
			Str("body", string(body)).
			Msg("Failed to create order")
		return fmt.Errorf("failed to create order: %s", string(body))
	}

	log.Info().Msg("Order created successfully")
	return nil
}

func (p *MockProvider) CreateOrder(ctx context.Context, order models.Order) error {

	log.Info().Msg("starting to create Wolt order")

	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Error().Err(err).Msg("error marshaling order to json")
		return err
	}

	req, err := http.NewRequest("POST", p.endpoint+"/v1/orders", bytes.NewBuffer(orderJSON))
	if err != nil {
		log.Error().Err(err).Msg("error creating http request")
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Error().Err(err).Msg("error making http request")
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("failed to read response body")
			return fmt.Errorf("failed to read response body: %w", err)
		}
		log.Error().
			Int("status", resp.StatusCode).
			Str("body", string(body)).
			Msg("Failed to create order")
		return fmt.Errorf("failed to create order: %s", string(body))
	}

	log.Info().Msg("Order created successfully")
	return nil
}
