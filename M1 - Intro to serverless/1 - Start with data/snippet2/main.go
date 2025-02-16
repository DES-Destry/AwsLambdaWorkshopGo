package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// LineItem represents a mock line item for an order.
type LineItem struct {
	SKU      int    `json:"sku"`
	Color    string `json:"color"`
	Quantity int    `json:"quantity"`
	InStock  bool   `json:"in_stock"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Create a mock line item for an order
	lineItem := LineItem{
		SKU:      1234242,
		Color:    "blue",
		Quantity: 42,
		InStock:  true,
	}

	body, err := json.Marshal(lineItem)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	lambda.Start(handler)
}
