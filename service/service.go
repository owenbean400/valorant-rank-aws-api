package service

import (
	"encoding/json"
	"errors"
	"net/http"
	valorantdao "valorant-rank-api/dao"
	"valorant-rank-api/domain/environment"
	"valorant-rank-api/domain/helper"
	"valorant-rank-api/domain/structure"
)

func UpdateDataWithAPI(puuid string) error {
	url := "https://api.henrikdev.xyz/valorant/v1/mmr-history/na" + puuid

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.New("error creating request: " + err.Error())
	}

	req.Header.Add("Authorization", environment.GetValorantAPIKeyEnv())
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("error making request: " + err.Error())
	}
	defer resp.Body.Close()

	var result structure.MrrData
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&result)

	if err != nil {
		return errors.New("error parsing JSON: " + err.Error())
	}

	for _, mrr_stats := range result.MrrStats {
		isMatchExist, err := valorantdao.MatchAlreadyExist(puuid, mrr_stats.MatchId)

		if err == nil && !isMatchExist {
			matchData, err := getMatchData(mrr_stats.MatchId)

			if err == nil {
				player := helper.GetPlayerByPUIID(matchData.Data.PlayersInfo.Players, puuid)

				if player != nil {
					var teamStat structure.TeamStat

					if player.Team == "Red" {
						teamStat = matchData.Data.Team.Red
					}
					if player.Team == "Blue" {
						teamStat = matchData.Data.Team.Blue
					}

					playMetaStatSave := structure.PlayerMetaStatSave{Score: player.PlayerStats.Score,
						Kills:          player.PlayerStats.Kills,
						Deaths:         player.PlayerStats.Deaths,
						Assists:        player.PlayerStats.Assists,
						BodyShots:      player.PlayerStats.BodyShots,
						HeadShots:      player.PlayerStats.HeadShots,
						LegShots:       player.PlayerStats.LegShot,
						DamageMade:     player.DamageMade,
						DamageRecieved: player.DamageReceived,
					}
					rankSaveGame := structure.RankStatGameSave{PUUID: puuid,
						MatchId:        mrr_stats.MatchId,
						MmrChange:      mrr_stats.MmmrChange,
						Map:            matchData.Data.MetaData.Map,
						Character:      player.Character,
						RoundsWon:      teamStat.RoundsWon,
						RoundsLost:     teamStat.RoundsLost,
						PlayerMetaStat: playMetaStatSave,
					}

					valorantdao.WriteRankValorantMatch(rankSaveGame)

				}
			}
		}
	}

	return nil
}

func GetValorantRankHistory(puuid string) ([]structure.RankStatGameSave, error) {
	var rank_games []structure.RankStatGameSave

	rank_games, err := valorantdao.QueryValorantMatches(puuid)

	if err != nil {
		return rank_games, errors.New("error getting Valorant rank history: " + err.Error())
	}

	return rank_games, nil
}

func getMatchData(match_id string) (structure.MatchData, error) {
	var result structure.MatchData

	url := "https://api.henrikdev.xyz/valorant/v2/match/" + match_id

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, errors.New("could not start new request: " + err.Error())
	}

	req.Header.Add("Authorization", environment.GetValorantAPIKeyEnv())
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, errors.New("error making request: " + err.Error())
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&result)

	if err != nil {
		return result, errors.New("error parsing json: " + err.Error())
	}

	return result, nil
}
