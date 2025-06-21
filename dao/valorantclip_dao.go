package dao

import (
	"errors"
	"fmt"
	"valorant-rank-api/domain/environment"
	"valorant-rank-api/domain/structure"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

func GetValorantClip(uuid string) (structure.ValorantClipJSON, error) {
	var valorant_clip_dynamodb structure.ValorantClipDynamoDbRecord
	var valorant_clip structure.ValorantClipJSON

	ctx, svc, err := GetDynamoDb()

	if err != nil {
		return valorant_clip, fmt.Errorf("error setting up Dynamo DB: %w", err)
	}

	tableName := environment.GetClipTableName()

	input := dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"uuid": &types.AttributeValueMemberS{Value: uuid},
		},
	}

	res, err := svc.GetItem(ctx, &input)

	if err != nil {
		return valorant_clip, fmt.Errorf("clip not found from DynamoDB: %w", err)
	}

	err = attributevalue.UnmarshalMap(res.Item, &valorant_clip_dynamodb)

	if err != nil {
		return valorant_clip, fmt.Errorf("error parsing dynamoDB data of clips: %w", err)
	}

	valorant_clip = structure.ValorantClipJSON{
		ID:         valorant_clip_dynamodb.ID,
		BaseURL:    valorant_clip_dynamodb.BaseURL,
		FileName:   valorant_clip_dynamodb.FileName,
		Extenstion: valorant_clip_dynamodb.Extenstion,
		FilePath:   valorant_clip_dynamodb.FilePath,
		FullUrl:    valorant_clip_dynamodb.FullUrl,
		MatchId:    valorant_clip_dynamodb.MatchId,
	}

	return valorant_clip, nil
}

func ScanValorantClips(queryAmount int32) ([]structure.ValorantClipJSON, error) {
	var valorant_clips_dynamodb []structure.ValorantClipDynamoDbRecord
	var valorant_clips []structure.ValorantClipJSON

	ctx, svc, err := GetDynamoDb()

	if err != nil {
		return valorant_clips, fmt.Errorf("error setting up Dynamo DB: %w", err)
	}

	tableName := environment.GetClipTableName()

	input := dynamodb.ScanInput{
		TableName: aws.String(tableName),
		Limit:     aws.Int32(queryAmount),
	}

	res, err := svc.Scan(ctx, &input)

	if err != nil {
		return valorant_clips, fmt.Errorf("error scan Dynamo DB: %w", err)
	}

	err = attributevalue.UnmarshalListOfMaps(res.Items, &valorant_clips_dynamodb)

	if err != nil {
		return valorant_clips, fmt.Errorf("unmarshal failed: %w", err)
	}

	for _, element := range valorant_clips_dynamodb {
		valorant_clips = append(valorant_clips, structure.ValorantClipJSON{
			ID:         element.ID,
			BaseURL:    element.BaseURL,
			FileName:   element.FileName,
			Extenstion: element.Extenstion,
			FilePath:   element.FilePath,
			FullUrl:    element.FullUrl,
			MatchId:    element.MatchId,
		})
	}

	return valorant_clips, nil
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
