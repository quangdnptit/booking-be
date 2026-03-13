package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// NewDynamoDBStore builds a DynamoDB client.
//
// Local DynamoDB (Docker / DynamoDB Local):
//   - Set DYNAMODB_ENDPOINT=http://localhost:8000
//   - The SDK uses static credentials (fake/fake) so broken ~/.aws or SSO sessions
//     do not affect local calls. Region defaults to us-east-1 unless AWS_REGION is set.
//
// Real AWS:
//   - Leave DYNAMODB_ENDPOINT empty
//   - Set valid AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_REGION
//   - Or use ~/.aws/credentials + AWS_PROFILE (must be valid, not expired SSO)
func NewDynamoDBStore(ctx context.Context, endpoint string) (*dynamodb.Client, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	var cfg aws.Config
	var err error

	if endpoint != "" {
		// Local: never use shared config credentials (often invalid / expired SSO)
		accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
		secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
		if accessKey == "" {
			accessKey = "fake"
		}
		if secretKey == "" {
			secretKey = "fake"
		}
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	}
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	var opts []func(*dynamodb.Options)
	if endpoint != "" {
		opts = append(opts, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}

	return dynamodb.NewFromConfig(cfg, opts...), nil
}
