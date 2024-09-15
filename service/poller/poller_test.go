package poller

import (
	"fmt"
	"github.com/heloayer/check-order-status/core/external_api/providers"
	"github.com/heloayer/check-order-status/internal/models"
	mocks "github.com/heloayer/check-order-status/mocks/core/external_api/providers"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestPoller_Success(t *testing.T) {

	mockProvider := mocks.NewProvider(t)

	orders := []models.Order{
		{ID: "1", Status: "delivered"},
		{ID: "2", Status: "new"},
	}

	mockProvider.On("GetOrders", mock.Anything).Return(orders, nil)

	poller := NewPoller([]providers.Provider{mockProvider}, 100*time.Millisecond)

	go poller.Start()
	defer poller.Stop()

	time.Sleep(200 * time.Millisecond)

	mockProvider.AssertExpectations(t)
}

func TestPoller_Failure(t *testing.T) {
	mockProvider := mocks.NewProvider(t)

	mockProvider.On("GetOrders", mock.Anything).Return(nil, fmt.Errorf("failed to get orders"))

	poller := NewPoller([]providers.Provider{mockProvider}, 100*time.Millisecond)

	go poller.Start()
	defer poller.Stop()

	time.Sleep(200 * time.Millisecond)

	mockProvider.AssertExpectations(t)
}
