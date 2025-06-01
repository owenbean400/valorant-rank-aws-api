package environment

import "os"

func GetPlayerPuuidEnv() string {
	return os.Getenv("PLAYER_PUUID")
}

func GetValorantAPIKeyEnv() string {
	return os.Getenv("VALORANT_API_KEY")
}

func GetRoleArn() string {
	return os.Getenv("AWS_DYNAMODB_ROLE_ARN")
}

func GetTableName() string {
	return os.Getenv("AWS_DYNAMODB_TABLE_ARN")
}

func GetSessionName() string {
	sessionNameOs := os.Getenv("AWS_DYNAMODB_SESSION_NAME")
	if sessionNameOs == "" {
		return "valorant-dynamodb-session"
	} else {
		return sessionNameOs
	}
}
