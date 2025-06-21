package dao

import (
	"context"
	"errors"
	"valorant-rank-api/domain/environment"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func GetDynamoDb() (context.Context, *dynamodb.Client, error) {
	var ctx context.Context
	var svc *dynamodb.Client

	ctx = context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return ctx, svc, errors.New("unable to load SDK config" + err.Error())
	}

	roleArn := environment.GetRoleArn()
	sessionName := environment.GetSessionName()

	stsClient := sts.NewFromConfig(cfg)

	creds := stscreds.NewAssumeRoleProvider(stsClient, roleArn, func(o *stscreds.AssumeRoleOptions) {
		o.RoleSessionName = sessionName
	})

	roleCfg := cfg.Copy()
	roleCfg.Credentials = aws.NewCredentialsCache(creds)

	svc = dynamodb.NewFromConfig(roleCfg)

	return ctx, svc, nil
}
