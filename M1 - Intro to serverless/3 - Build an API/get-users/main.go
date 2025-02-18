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

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (resp events.APIGatewayProxyResponse, err error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	resp = events.APIGatewayProxyResponse{StatusCode: 500}
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return
	}

	client := dynamodb.NewFromConfig(cfg)

	input := &dynamodb.ScanInput{
		TableName: aws.String("serverless_workshop_intro"),
	}

	result, err := client.Scan(ctx, input)
	if err != nil {
		log.Printf("Error scanning table: %v", err)
		return
	}

	var items []map[string]interface{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		log.Printf("Error unmarshalling items: %v", err)
		return
	}

	body, err := json.Marshal(items)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		return
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
