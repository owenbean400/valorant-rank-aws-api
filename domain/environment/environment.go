package environment

import "os"

func GetPlayerPuuidEnv() string {
	return os.Getenv("PLAYER_PUUID")
}

func GetValorantAPIKeyEnv() string {
	return os.Getenv("VALORANT_API_KEY")
}

func GetRankTableName() string {
	return os.Getenv("AWS_DYNAMODB_RANK_TABLE_ARN")
}

func GetClipTableName() string {
	return os.Getenv("AWS_DYNAMODB_CLIP_TABLE_ARN")
}

func GetSessionName() string {
	sessionNameOs := os.Getenv("AWS_DYNAMODB_SESSION_NAME")
	if sessionNameOs == "" {
		return "valorant-dynamodb-session"
	} else {
		return sessionNameOs
	}
}
