package models

type Order struct {
	ID          string    `json:"id"`
	Status      string    `json:"status,omitempty"`
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
