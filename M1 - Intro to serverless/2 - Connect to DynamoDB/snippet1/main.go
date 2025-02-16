package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tableName := "serverless_workshop_intro"

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("failed to load config: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	// Create DynamoDB client
	client := dynamodb.NewFromConfig(cfg)

	// Define the list of people
	people := []struct {
		Userid string
		Name   string
	}{
		{"marivera", "Martha Rivera"},
		{"nikkwolf", "Nikki Wolf"},
		{"pasantos", "Paulo Santos"},
	}

	// Prepare write requests for batch write
	writeRequests := make([]types.WriteRequest, 0, len(people))
	for _, person := range people {
		// Generate a unique ID (UUID without hyphens)
		id := uuid.New().String() // To remove hyphens: strings.ReplaceAll(uuid.New().String(), "-", "")
		item := map[string]types.AttributeValue{
			"_id":      &types.AttributeValueMemberS{Value: id},
			"Userid":   &types.AttributeValueMemberS{Value: person.Userid},
			"FullName": &types.AttributeValueMemberS{Value: person.Name},
		}

		log.Printf("> batch writing: %s", person.Userid)

		writeRequests = append(writeRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: item,
			},
		})
	}

	// Create the input for BatchWriteItem
	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			tableName: writeRequests,
		},
	}

	_, err = client.BatchWriteItem(ctx, input)
	if err != nil {
		log.Printf("failed to batch write items: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	result := fmt.Sprintf("Success. Added %d people to %s.", len(people), tableName)
	responseBody, err := json.Marshal(map[string]string{"message": result})
	if err != nil {
		log.Printf("failed to marshal response: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	lambda.Start(handler)
}
