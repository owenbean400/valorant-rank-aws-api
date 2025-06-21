package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"valorant-rank-api/dao"
	"valorant-rank-api/domain/structure"
)

func GetValorantClips(body string) ([]structure.ValorantClipJSON, error, int) {
	var valorant_clips []structure.ValorantClipJSON
	var query_data structure.ValorantClipQuery
	var default_page int32 = 10

	err := json.Unmarshal([]byte(body), &query_data)

	if err != nil {
		return valorant_clips, fmt.Errorf("error parsing body: %w", err), 403
	}

	if query_data.PageLength == nil || *query_data.PageLength > 10 {
		query_data.PageLength = &default_page
	}

	if *query_data.PageLength > 10 || *query_data.PageLength < 1 {
		return valorant_clips, fmt.Errorf("integer is outside range of 1-10 query page"), 403
	}

	valorant_clips, err = dao.ScanValorantClips(*query_data.PageLength)

	if err != nil {
		return valorant_clips, fmt.Errorf("error scanning clips: %w", err), 500
	}

	return valorant_clips, nil, 200
}

func GetValorantClip(uuid string) (structure.ValorantClipJSON, error, int) {
	var valorant_clip structure.ValorantClipJSON

	if uuid == "" {
		return valorant_clip, nil, 404
	}

	valorant_clip, err := dao.GetValorantClip(uuid)

	if err != nil {
		if strings.HasPrefix(err.Error(), "clip not found from DynamoDB") {
			return valorant_clip, err, 404
		} else {
			return valorant_clip, err, 500
		}
	}

	return valorant_clip, nil, 302
}
