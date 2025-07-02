package service

import (
	"fmt"
	"strconv"
	"strings"
	"valorant-rank-api/dao"
	"valorant-rank-api/domain/structure"
)

func GetValorantClips(page_number_str string, last_eval_key string) (structure.ValorantClipsTable, int, error) {
	var valorant_clips_table structure.ValorantClipsTable

	page_number, err := strconv.ParseInt(page_number_str, 10, 32)
	if err != nil {
		return valorant_clips_table, 403, fmt.Errorf("error parsing parameter query pageLength as not an integer: %w", err)
	}

	if page_number > 10 || page_number < 1 {
		return valorant_clips_table, 403, fmt.Errorf("parameter query pageLength integer is outside range of 1-10 query page")
	}

	valorant_clips_table, err = dao.ScanValorantClips(int32(page_number), last_eval_key)

	if err != nil {
		return valorant_clips_table, 500, fmt.Errorf("error scanning clips: %w", err)
	}

	return valorant_clips_table, 200, nil
}

func GetValorantClip(uuid string) (structure.ValorantClipJSON, int, error) {
	var valorant_clip structure.ValorantClipJSON

	if uuid == "" {
		return valorant_clip, 404, nil
	}

	valorant_clip, err := dao.GetValorantClip(uuid)

	if err != nil {
		if strings.HasPrefix(err.Error(), "clip not found from DynamoDB") {
			return valorant_clip, 404, err
		} else {
			return valorant_clip, 500, err
		}
	}

	return valorant_clip, 302, nil
}
