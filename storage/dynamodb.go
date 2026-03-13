package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamo "github.com/guregu/dynamo/v2"
)

// NewDynamoDBStore builds AWS config and a guregu/dynamo DB
// Local: DYNAMODB_ENDPOINT=http://localhost:8000, static fake credentials.
// AWS: leave DYNAMODB_ENDPOINT empty, use normal credentials.
func NewDynamoDBStore(ctx context.Context, endpoint string) (*dynamo.DB, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	var cfg aws.Config
	var err error

	if endpoint != "" {
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
			config.WithRetryer(func() aws.Retryer {
				return retry.NewStandard(dynamo.RetryTxConflicts)
			}),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithRetryer(func() aws.Retryer {
				return retry.NewStandard(dynamo.RetryTxConflicts)
			}),
		)
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
	return dynamo.New(cfg, opts...), nil
}
