package dao

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"valorant-rank-api/domain/environment"
	"valorant-rank-api/domain/structure"
)

func WriteRankValorantMatch(valorant_data structure.RankStatGameSave) error {
	puuid_match_record := structure.ValorantRankDynamoDbRecord{
		KeyPuuidMatch:  getPrimaryKey(valorant_data.PUUID, &valorant_data.MatchId),
		RawDateInt:     valorant_data.RawDateInt,
		PUUID:          valorant_data.PUUID,
		MatchId:        valorant_data.MatchId,
		MmrChange:      valorant_data.MmrChange,
		DateStr:        valorant_data.DateStr,
		Map:            valorant_data.Map,
		Character:      valorant_data.Character,
		RoundsWon:      valorant_data.RoundsWon,
		RoundsLost:     valorant_data.RoundsLost,
		Score:          valorant_data.PlayerMetaStat.Score,
		Kills:          valorant_data.PlayerMetaStat.Kills,
		Deaths:         valorant_data.PlayerMetaStat.Deaths,
		Assists:        valorant_data.PlayerMetaStat.Assists,
		BodyShots:      valorant_data.PlayerMetaStat.BodyShots,
		HeadShots:      valorant_data.PlayerMetaStat.HeadShots,
		LegShots:       valorant_data.PlayerMetaStat.LegShots,
		DamageMade:     valorant_data.PlayerMetaStat.DamageMade,
		DamageRecieved: valorant_data.PlayerMetaStat.DamageRecieved,
	}

	puuid_record := structure.ValorantRankDynamoDbRecord{
		KeyPuuidMatch:  getPrimaryKey(valorant_data.PUUID, nil),
		RawDateInt:     valorant_data.RawDateInt,
		PUUID:          valorant_data.PUUID,
		MatchId:        valorant_data.MatchId,
		MmrChange:      valorant_data.MmrChange,
		DateStr:        valorant_data.DateStr,
		Map:            valorant_data.Map,
		Character:      valorant_data.Character,
		RoundsWon:      valorant_data.RoundsWon,
		RoundsLost:     valorant_data.RoundsLost,
		Score:          valorant_data.PlayerMetaStat.Score,
		Kills:          valorant_data.PlayerMetaStat.Kills,
		Deaths:         valorant_data.PlayerMetaStat.Deaths,
		Assists:        valorant_data.PlayerMetaStat.Assists,
		BodyShots:      valorant_data.PlayerMetaStat.BodyShots,
		HeadShots:      valorant_data.PlayerMetaStat.HeadShots,
		LegShots:       valorant_data.PlayerMetaStat.LegShots,
		DamageMade:     valorant_data.PlayerMetaStat.DamageMade,
		DamageRecieved: valorant_data.PlayerMetaStat.DamageRecieved,
	}

	err := saveValorantTable(puuid_match_record)

	if err != nil {
		return fmt.Errorf("error saving Valorant Rank Match ID: %w", err)
	}

	err = saveValorantTable(puuid_record)

	if err != nil {
		return fmt.Errorf("error saving Valorant Rank Player Change: %w", err)
	}

	return nil
}

func QueryValorantMatches(puuid string, pageNumber int32, start_eval_key_puuid_match string) (structure.ValorantRankHistoryTable, error) {
	var valorant_matches_dynamodb []structure.ValorantRankDynamoDbRecord
	var valorant_matches []structure.RankStatGameSave
	var valorant_matches_json structure.ValorantRankHistoryTable

	ctx, svc, err := GetDynamoDb()

	if err != nil {
		return valorant_matches_json, fmt.Errorf("error setting up Dynamo DB: %w", err)
	}

	tableName := environment.GetRankTableName()

	var input dynamodb.QueryInput
	if start_eval_key_puuid_match == "" {
		input = dynamodb.QueryInput{
			TableName:              aws.String(tableName),
			KeyConditionExpression: aws.String("puuid_match = :puuid_match"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":puuid_match": &types.AttributeValueMemberS{Value: puuid},
			},
			ScanIndexForward: aws.Bool(false),
			Limit:            aws.Int32(pageNumber),
		}
	} else {
		input = dynamodb.QueryInput{
			TableName:              aws.String(tableName),
			KeyConditionExpression: aws.String("puuid_match = :puuid_match"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":puuid_match": &types.AttributeValueMemberS{Value: puuid},
			},
			ScanIndexForward:  aws.Bool(false),
			Limit:             aws.Int32(pageNumber),
			ExclusiveStartKey: map[string]types.AttributeValue{"puuid_match": &types.AttributeValueMemberS{Value: start_eval_key_puuid_match}},
		}
	}

	res, err := svc.Query(ctx, &input)

	if err != nil {
		return valorant_matches_json, fmt.Errorf("error query up Dynamo DB: %w", err)
	}

	var last_eval_key_puuid_match string

	if val, ok := res.LastEvaluatedKey["puuid_match"]; ok {
		if s, ok := val.(*types.AttributeValueMemberS); ok {
			last_eval_key_puuid_match = s.Value
		} else {
			last_eval_key_puuid_match = ""
		}
	} else {
		last_eval_key_puuid_match = ""
	}

	err = attributevalue.UnmarshalListOfMaps(res.Items, &valorant_matches_dynamodb)
	if err != nil {
		return valorant_matches_json, fmt.Errorf("unmarshal failed: %w", err)
	}

	for _, element := range valorant_matches_dynamodb {
		player_meta_stat_save := structure.PlayerMetaStatSave{
			Score:          element.Score,
			Kills:          element.Kills,
			Deaths:         element.Deaths,
			Assists:        element.Assists,
			BodyShots:      element.BodyShots,
			HeadShots:      element.HeadShots,
			LegShots:       element.LegShots,
			DamageMade:     element.DamageMade,
			DamageRecieved: element.DamageRecieved,
		}

		valorant_matches = append(valorant_matches, structure.RankStatGameSave{
			PUUID:          element.PUUID,
			MatchId:        element.MatchId,
			RawDateInt:     element.RawDateInt,
			DateStr:        element.DateStr,
			MmrChange:      element.MmrChange,
			Map:            element.Map,
			Character:      element.Character,
			PlayerMetaStat: player_meta_stat_save,
			RoundsWon:      element.RoundsWon,
			RoundsLost:     element.RoundsLost,
		})
	}

	return structure.ValorantRankHistoryTable{
		History:                    valorant_matches,
		LastEvaluatedKeyPuuidMatch: last_eval_key_puuid_match,
	}, nil
}

func DoesMatchExist(puuid string, match_id string, rawDateInt int) (bool, error) {
	ctx, svc, err := GetDynamoDb()

	if err != nil {
		return false, fmt.Errorf("error setting up DynamoDB: %w", err)
	}

	tableName := environment.GetRankTableName()

	key := map[string]types.AttributeValue{
		"puuid_match":  &types.AttributeValueMemberS{Value: getPrimaryKey(puuid, &match_id)},
		"raw_date_int": &types.AttributeValueMemberN{Value: strconv.Itoa(rawDateInt)},
	}

	result, err := svc.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	})

	if err != nil {
		return false, fmt.Errorf("failed to get item: %w", err)
	}

	if result.Item == nil {
		return false, nil
	} else {
		return true, nil
	}
}

func saveValorantTable(valorant_rank_item structure.ValorantRankDynamoDbRecord) error {
	ctx, svc, err := GetDynamoDb()

	if err != nil {
		return fmt.Errorf("error setting up DynamoDB: %w", err)
	}

	tableName := environment.GetRankTableName()

	av, err := attributevalue.MarshalMap(valorant_rank_item)
	if err != nil {
		return fmt.Errorf("error parsing Valorant Rank data: %w", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(ctx, input)

	if err != nil {
		return fmt.Errorf("error putItem in Valorant Rank DynamoDB: %w", err)
	}

	return nil
}

func getPrimaryKey(puuid string, match_id *string) string {
	if match_id == nil {
		return puuid
	} else {
		return puuid + "::" + *match_id
	}
}
