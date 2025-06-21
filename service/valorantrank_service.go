package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	valorantdao "valorant-rank-api/dao"
	"valorant-rank-api/domain/environment"
	"valorant-rank-api/domain/helper"
	"valorant-rank-api/domain/structure"
)

func UpdateDataWithAPI(puuid string) error {
	errorsStr := ""

	url := "https://api.henrikdev.xyz/valorant/v1/by-puuid/mmr-history/na/" + puuid

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", environment.GetValorantAPIKeyEnv())
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	defer resp.Body.Close()

	var result structure.MrrData
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&result)

	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	for _, mrr_stats := range result.MrrStats {
		isMatchExist, err := valorantdao.DoesMatchExist(puuid, mrr_stats.MatchId, mrr_stats.DateRaw)

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
						RawDateInt:     mrr_stats.DateRaw,
						DateStr:        mrr_stats.Date,
						Map:            matchData.Data.MetaData.Map,
						Character:      player.Character,
						RoundsWon:      teamStat.RoundsWon,
						RoundsLost:     teamStat.RoundsLost,
						PlayerMetaStat: playMetaStatSave,
					}

					err = valorantdao.WriteRankValorantMatch(rankSaveGame)

					if err != nil {
						errorsStr += err.Error() + "::"
					}

				} else {
					errorsStr += "error with finding player ID ::"
				}
			} else {
				errorsStr += err.Error() + "::"
			}
		} else if err != nil {
			errorsStr += err.Error() + "::"
		}
	}

	if errorsStr != "" {
		return errors.New(errorsStr)
	}

	return nil
}

func GetValorantRankHistory(puuid string) ([]structure.RankStatGameSave, error) {
	var rank_games []structure.RankStatGameSave

	rank_games, err := valorantdao.QueryValorantMatches(puuid)

	if err != nil {
		return rank_games, fmt.Errorf("error getting Valorant rank history: %w", err)
	}

	return rank_games, nil
}

func getMatchData(match_id string) (structure.MatchData, error) {
	var result structure.MatchData

	url := "https://api.henrikdev.xyz/valorant/v2/match/" + match_id

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, fmt.Errorf("could not start new request: %w", err)
	}

	req.Header.Add("Authorization", environment.GetValorantAPIKeyEnv())
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&result)

	if err != nil {
		return result, fmt.Errorf("error parsing json: %w", err)
	}

	return result, nil
}
