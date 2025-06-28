package service

import (
	"fmt"
	"strconv"
	"strings"
	"valorant-rank-api/dao"
	"valorant-rank-api/domain/structure"
)

func GetValorantClips(pageNumberStr string) ([]structure.ValorantClipJSON, error, int) {
	var valorant_clips []structure.ValorantClipJSON

	pageNumber, err := strconv.ParseInt(pageNumberStr, 10, 32)
	if err != nil {
		return valorant_clips, fmt.Errorf("error parsing parameter query pageLength as not an integer: %w", err), 403
	}

	if pageNumber > 10 || pageNumber < 1 {
		return valorant_clips, fmt.Errorf("parameter query pageLength integer is outside range of 1-10 query page"), 403
	}

	valorant_clips, err = dao.ScanValorantClips(int32(pageNumber))

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
