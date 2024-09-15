package order

import (
	"context"
	"errors"
	"github.com/heloayer/check-order-status/internal/models"
	"github.com/heloayer/check-order-status/internal/resources/http/handler/dto"
	mocks "github.com/heloayer/check-order-status/mocks/internal_/service/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateOrder(t *testing.T) {

	mockOrder := mocks.NewOrder(t)

	request := dto.OrderRequest{
		Customer: dto.Customer{
			Name:        "Gordon Ramsay",
			PhoneNumber: "+78017584677",
		},
		TotalSum:    172.5,
		PaymentType: "card",
		Products: []dto.Product{
			{
				Name:     "royal spaghetti",
				Quantity: 1,
				Price:    12.5,
				Comment:  "",
			},
			{
				Name:     "wine 1966",
				Quantity: 1,
				Price:    160,
				Comment:  "",
			},
		},
	}

	mockOrder.On("CreateOrder", mock.Anything, request, "glovo").Return(nil)

	mockOrder.On("CreateOrder", mock.Anything, request, "wolt").Return(nil)

	err := mockOrder.CreateOrder(context.Background(), request, "glovo")
	assert.NoError(t, err)

	err = mockOrder.CreateOrder(context.Background(), request, "wolt")
	assert.NoError(t, err)

	mockOrder.AssertExpectations(t)
}

func TestCreateOrder_Failure(t *testing.T) {

	mockOrder := mocks.NewOrder(t)

	request := dto.OrderRequest{
		Customer: dto.Customer{
			Name:        "Gordon Ramsay",
			PhoneNumber: "+78017584677",
		},
		TotalSum:    172.5,
		PaymentType: "card",
		Products: []dto.Product{
			{
				Name:     "royal spaghetti",
				Quantity: 1,
				Price:    12.5,
				Comment:  "",
			},
			{
				Name:     "wine 1966",
				Quantity: 1,
				Price:    160,
				Comment:  "",
			},
		},
	}

	mockOrder.On("CreateOrder", mock.Anything, request, "glovo").Return(errors.New("failed to create order"))

	err := mockOrder.CreateOrder(context.Background(), request, "glovo")
	assert.Error(t, err)
	assert.Equal(t, "failed to create order", err.Error())

	mockOrder.AssertExpectations(t)
}

func TestGetOrderByID(t *testing.T) {

	mockOrder := mocks.NewOrder(t)

	expectedOrder := models.Order{
		ID:     "8",
		Status: "new",
		Customer: models.Customer{
			Name:        "Michael",
			PhoneNumber: "+98006573487",
		},
		TotalSum:    23.5,
		PaymentType: "card",
		Products: []models.Product{
			{
				Name:     "Bowl of meat",
				Quantity: 1,
				Price:    8.5,
				Comment:  "chicken",
			},
			{
				Name:     "champagne",
				Quantity: 1,
				Price:    15,
				Comment:  "sparkle",
			},
		},
	}

	mockOrder.On("GetOrderByID", mock.Anything, "1").Return(expectedOrder, nil)

	order, err := mockOrder.GetOrderByID(context.Background(), "1")
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)

	mockOrder.AssertExpectations(t)
}

func TestGetOrderByID_Failure(t *testing.T) {
	mockOrder := mocks.NewOrder(t)

	mockOrder.On("GetOrderByID", mock.Anything, "2").Return(models.Order{}, errors.New("order not found"))

	order, err := mockOrder.GetOrderByID(context.Background(), "2")
	assert.Error(t, err)
	assert.Equal(t, "order not found", err.Error())
	assert.Equal(t, models.Order{}, order)

	mockOrder.AssertExpectations(t)
}
