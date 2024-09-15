package dto

import "github.com/heloayer/check-order-status/internal/models"

type OrderRequest struct {
	Customer    Customer  `json:"customer"`
	TotalSum    float64   `json:"total_sum"`
	PaymentType string    `json:"payment_type"`
	Products    []Product `json:"products"`
}

type Customer struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Comment  string  `json:"comment,omitempty"`
}

func (or OrderRequest) ToModel() models.Order {
	order := models.Order{
		Customer: models.Customer{
			Name:        or.Customer.Name,
			PhoneNumber: or.Customer.PhoneNumber,
		},
		TotalSum:    or.TotalSum,
		PaymentType: or.PaymentType,
	}

	for _, product := range or.Products {
		order.Products = append(order.Products, models.Product{
			Name:     product.Name,
			Quantity: product.Quantity,
			Price:    product.Price,
			Comment:  product.Comment,
		})
	}
	order.Status = "new"
	return order
}
