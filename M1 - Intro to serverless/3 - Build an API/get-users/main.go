package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	client := dynamodb.NewFromConfig(cfg)

	input := &dynamodb.ScanInput{
		TableName: aws.String("serverless_workshop_intro"),
	}

	result, err := client.Scan(ctx, input)
	if err != nil {
		log.Printf("Error scanning table: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	var items []map[string]interface{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		log.Printf("Error unmarshalling items: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	body, err := json.Marshal(items)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
