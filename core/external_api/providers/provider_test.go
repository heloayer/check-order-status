package providers

import (
	"context"
	"fmt"
	"testing"

	"github.com/heloayer/check-order-status/internal/models"
	"github.com/heloayer/check-order-status/mocks/core/external_api/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrder(t *testing.T) {
	mockProvider := mocks.NewProvider(t)

	order := models.Order{
		ID:     "1",
		Status: "new",
		Customer: models.Customer{
			Name:        "Billy",
			PhoneNumber: "+78007684376",
		},
		TotalSum:    7.9,
		PaymentType: "card",
		Products: []models.Product{
			{
				Name:     "zero coca cola",
				Quantity: 1,
				Price:    2.5,
				Comment:  "Zero sugar",
			},
			{
				Name:     "pizza",
				Quantity: 1,
				Price:    5.4,
				Comment:  "Italian pizza with a lot of cheese",
			},
		},
	}

	mockProvider.On("CreateOrder", mock.Anything, order).Return(nil)

	err := mockProvider.CreateOrder(context.Background(), order)
	assert.NoError(t, err)

	mockProvider.AssertExpectations(t)
}

func TestCreateOrder_Failure(t *testing.T) {

	mockProvider := mocks.NewProvider(t)

	order := models.Order{
		ID:     "1",
		Status: "new",
		Customer: models.Customer{
			Name:        "Billy",
			PhoneNumber: "+78007684376",
		},
		TotalSum:    7.9,
		PaymentType: "card",
		Products: []models.Product{
			{
				Name:     "zero coca cola",
				Quantity: 1,
				Price:    2.5,
				Comment:  "Zero sugar",
			},
			{
				Name:     "pizza",
				Quantity: 1,
				Price:    5.4,
				Comment:  "Italian pizza with a lot of cheese",
			},
		},
	}

	mockProvider.On("CreateOrder", mock.Anything, order).Return(fmt.Errorf("failed to create order"))

	err := mockProvider.CreateOrder(context.Background(), order)
	assert.Error(t, err)
	assert.Equal(t, "failed to create order", err.Error())

	mockProvider.AssertExpectations(t)
}

func TestGetOrders(t *testing.T) {

	mockProvider := mocks.NewProvider(t)

	orders := []models.Order{
		{ID: "1", Status: "delivered"},
		{ID: "2", Status: "new"},
		{ID: "3", Status: "ready_for_pickup"},
		{ID: "4", Status: "cancelled"},
		{ID: "5", Status: "cooking started"},
		{ID: "5", Status: "cooking complete"},
	}

	mockProvider.On("GetOrders", mock.Anything).Return(orders, nil)

	result, err := mockProvider.GetOrders(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, orders, result)

	mockProvider.AssertExpectations(t)
}

func TestGetOrders_Failure(t *testing.T) {
	mockProvider := mocks.NewProvider(t)

	mockProvider.On("GetOrders", mock.Anything).Return(nil, fmt.Errorf("failed to get orders"))

	result, err := mockProvider.GetOrders(context.Background())
	assert.Error(t, err)
	assert.Equal(t, "failed to get orders", err.Error())
	assert.Nil(t, result)

	mockProvider.AssertExpectations(t)
}
