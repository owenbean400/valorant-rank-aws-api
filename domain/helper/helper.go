package helper

import "valorant-rank-api/domain/structure"

func GetPlayerByPUIID(players []structure.Player, puuid string) *structure.Player {
	var result *structure.Player

	for i := range players {
		if players[i].PUUID == puuid {
			result = &players[i]
			break
		}
	}

	return result
}
