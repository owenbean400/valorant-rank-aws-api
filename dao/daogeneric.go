package dao

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func GetDynamoDb() (context.Context, *dynamodb.Client, error) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return ctx, nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	db := dynamodb.NewFromConfig(cfg)

	return ctx, db, nil
}
