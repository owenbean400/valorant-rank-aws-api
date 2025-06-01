package valorantdao

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"

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
		return errors.New("error saving Valorant Rank Match ID: " + err.Error())
	}

	err = saveValorantTable(puuid_record)

	if err != nil {
		return errors.New("error saving Valorant Rank Player Change: " + err.Error())
	}

	return nil
}

func QueryValorantMatches(puuid string) ([]structure.RankStatGameSave, error) {
	var valorant_matches_dynamodb []structure.ValorantRankDynamoDbRecord
	var valorant_matches []structure.RankStatGameSave

	ctx, svc, err := getDynamoDb()

	if err != nil {
		return valorant_matches, errors.New("error setting up Dynamo DB: " + err.Error())
	}

	tableName := environment.GetTableName()

	input := dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("puuid_match = :puuid_match"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":puuid_match": &types.AttributeValueMemberS{Value: puuid},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(10),
	}

	queryPaginator := dynamodb.NewQueryPaginator(svc, &input)

	for queryPaginator.HasMorePages() {
		response, err := queryPaginator.NextPage(ctx)
		if err != nil {
			break
		} else {
			var valorant_matches_dynamodb_page []structure.ValorantRankDynamoDbRecord
			err = attributevalue.UnmarshalListOfMaps(response.Items, &valorant_matches_dynamodb_page)
			if err != nil {
				break
			} else {
				valorant_matches_dynamodb = append(valorant_matches_dynamodb, valorant_matches_dynamodb_page...)
			}
		}
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

	return valorant_matches, nil
}

func MatchAlreadyExist(puuid string, match_id string) (bool, error) {
	ctx, svc, err := getDynamoDb()

	if err != nil {
		return false, errors.New("error with setup DynamoDB: " + err.Error())
	}

	tableName := environment.GetTableName()

	key := map[string]types.AttributeValue{
		"puuid_match": &types.AttributeValueMemberS{Value: getPrimaryKey(puuid, &match_id)},
	}

	result, err := svc.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	})

	if err != nil {
		log.Fatalf("failed to get item: %v", err)
	}

	if result.Item == nil {
		return true, nil
	} else {
		return false, nil
	}
}

func saveValorantTable(valorant_rank_item structure.ValorantRankDynamoDbRecord) error {
	ctx, svc, err := getDynamoDb()

	if err != nil {
		return errors.New("error setting up DynamoDB: " + err.Error())
	}

	tableName := environment.GetTableName()

	av, err := attributevalue.MarshalMap(valorant_rank_item)
	if err != nil {
		return errors.New("error parsing Valorant Rank data: " + err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(ctx, input)

	if err != nil {
		return errors.New("error putItem in Valorant Rank DynamoDB: " + err.Error())
	}

	return nil
}

func getDynamoDb() (context.Context, *dynamodb.Client, error) {
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

func getPrimaryKey(puuid string, match_id *string) string {
	if match_id == nil {
		return puuid
	} else {
		return puuid + "::" + *match_id
	}
}
