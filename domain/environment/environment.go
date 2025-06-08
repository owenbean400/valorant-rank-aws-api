package environment

import "os"

func GetPlayerPuuidEnv() string {
	return os.Getenv("PLAYER_PUUID")
}

func GetValorantAPIKeyEnv() string {
	return os.Getenv("VALORANT_API_KEY")
}

func GetTableName() string {
	return os.Getenv("AWS_DYNAMODB_TABLE_ARN")
}
