package dao

import (
	"errors"
	"valorant-rank-api/domain/environment"
	"valorant-rank-api/domain/structure"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func WriteClip(clip_data structure.ValorantClipJSON) error {
	clip_record := structure.ValorantClipDynamoDbRecord{
		ID:         clip_data.ID,
		BaseURL:    clip_data.BaseURL,
		FileName:   clip_data.FileName,
		Extenstion: clip_data.Extenstion,
		FilePath:   clip_data.FilePath,
		FullUrl:    clip_data.FullUrl,
	}

	err := saveClipTable(clip_record)

	if err != nil {
		return errors.New("error saving Valorant clip: " + err.Error())
	}

	return nil
}

func saveClipTable(clip structure.ValorantClipDynamoDbRecord) error {
	ctx, svc, err := GetDynamoDb()

	if err != nil {
		return errors.New("error setting up DynamoDB: " + err.Error())
	}

	tableName := environment.GetClipTableName()

	av, err := attributevalue.MarshalMap(clip)

	if err != nil {
		return errors.New("error parsing Valorant clip data: " + err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(ctx, input)

	if err != nil {
		return errors.New("error putItem in Valorant Clip DynamoDB: " + err.Error())
	}

	return nil
}
