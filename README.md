
# Check Order Status

Application for receiving order statuses from various providers. There is a functionality for periodically receiving order statuses and creating order or receiving full information of it.

Providers are JSON Server and Mock API.

## Installation

### JSON Server

Source: `https://github.com/typicode/json-server/tree/v0`

Choose one of the 3 options to install JSON Server:
- `npm install -g json-server@0.17.4`
- `npm install -g yarn`
- `yarn global add json-server@0.17.4`

### Mock API

No actions here. I used this public API's url in configuration.

## Endpoint Guide

### GET /v1/order/get-order/:order_id

Responses with body containing full order information

Mock API Contains up to 6 orders, IDs from 6 to 11.

JSON Server has 5 orders, IDs from 1 to 5.


### POST /v1/order/create-order

Requires: query and json body

- Query: Key: provider | Value: "wolt" or "glovo".
  Wolt uses MockAPI for creating, Glovo uses JSON Server.
- JSON boby

Sample Json Body:
`{ 
    "customer": {
    "name": "John Wayne",
    "phone_number": "+88007006545"
    },
    "total_sum": 11.99,
    "payment_type": "card",
    "products": [
        {
        "name": "cheeseburger",
        "quantity": 3,
        "price": 8.00,
        "comment": "Special offer"
        },
        {
        "name": "red wine",
        "quantity": 1,
        "price": 3.99,
        "comment": "Limited edition"
        }
      ]
}`

## Requirements

go version 1.23

## Run

1. Before running app, start JSON Server by executing:

- `json-server --watch db.json`

This command will read from db.json file placed in the main directory.

JSON Server uses 3000 port.

2. Run App

- `go run cmd/app/main.go`

App uses 8080 port

## Author

Ansar Issabekov

Linkedin: `https://www.linkedin.com/in/ansarissabekov/`

Telegram: `@cinephile0`

Email: `us.gogle11@gmail.com`

