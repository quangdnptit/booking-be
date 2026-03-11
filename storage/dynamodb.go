package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoDB implements Store interface with DynamoDB
type DynamoDB struct {
	client        *dynamodb.Client
	roomsTable    string
	bookingsTable string
}

// NewDynamoDBStore creates and initializes a DynamoDB store
func NewDynamoDBStore(ctx context.Context, endpoint string) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	// Allow custom endpoint for local development
	clientOptions := []func(*dynamodb.Options){}
	if endpoint != "" {
		clientOptions = append(clientOptions, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}

	return dynamodb.NewFromConfig(cfg, clientOptions...), nil
}
