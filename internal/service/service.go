package service

import "github.com/heloayer/check-order-status/internal/service/order"

type Service struct {
	Order order.OrderService
}

func New(order order.OrderService) *Service {
	return &Service{
		Order: order,
	}
}
