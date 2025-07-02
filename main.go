package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"valorant-rank-api/domain/environment"
	"valorant-rank-api/domain/version"
	"valorant-rank-api/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	switch req.RequestContext.HTTP.Method {
	case http.MethodGet:
		if strings.HasPrefix(req.RawPath, "/clip/") {
			uuid := req.RawPath[6:]

			clip, status_code, err := service.GetValorantClip(uuid)
			switch status_code {
			case 302:
				return redirectReponse(http.StatusFound, clip.FullUrl)
			case 404:
				return jsonResponse(http.StatusNotFound, map[string]string{"error": "not found"})
			default:
				return jsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		} else {
			switch req.RawPath {
			case "/clips":
				queryParams := req.QueryStringParameters

				pageNumberStr, ok := queryParams["pageLength"]
				if !ok {
					pageNumberStr = "10"
				}

				last_eval_key, ok := queryParams["lastEvalKey"]
				if !ok {
					last_eval_key = ""
				}

				clips_data, status_code, err := service.GetValorantClips(pageNumberStr, last_eval_key)

				switch status_code {
				case 200:
					return jsonResponse(http.StatusOK, clips_data)
				case 403:
					return jsonResponse(http.StatusBadRequest, map[string]string{"msg": "page_length must be an integer equal to or less than 10 and greater than 0"})
				default:
					return jsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()})
				}
			case "/history":
				query_params := req.QueryStringParameters

				page_number_str, ok := query_params["pageLength"]
				if !ok {
					page_number_str = "10"
				}

				last_eval_key_puuid_match, ok := query_params["lastEvalKeyPuuidMatch"]
				if !ok {
					last_eval_key_puuid_match = ""
				}

				last_eval_key_raw_date_int_str, ok := query_params["lastEvalKeyRawDateInt"]
				if !ok || last_eval_key_raw_date_int_str == "" {
					last_eval_key_raw_date_int_str = "-1"
				}

				rank_history, status_code, err := service.GetValorantRankHistory(environment.GetPlayerPuuidEnv(), page_number_str, last_eval_key_puuid_match, last_eval_key_raw_date_int_str)

				switch status_code {
				case 200:
					return jsonResponse(http.StatusOK, rank_history)
				case 403:
					return jsonResponse(http.StatusBadRequest, map[string]string{"msg": "page_length must be an integer equal to or less than 10 and greater than 0"})
				default:
					return jsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()})
				}
			case "/version":
				return jsonResponse(http.StatusOK, map[string]string{"version": version.GetVersionNumber()})

			default:
				return jsonResponse(http.StatusNotFound, map[string]string{"error": "not found"})
			}
		}
	case http.MethodPost:
		switch req.RawPath {
		case "/update":
			err := service.UpdateDataWithAPI(environment.GetPlayerPuuidEnv())

			if err != nil {
				return jsonResponse(http.StatusInternalServerError, map[string]string{"error": "internal error"})
			} else {
				return jsonResponse(http.StatusOK, map[string]string{"msg": "data updated"})
			}

		default:
			return jsonResponse(http.StatusNotFound, map[string]string{"error": "not found"})
		}
	default:
		return jsonResponse(http.StatusNotFound, map[string]string{"error": "not found"})
	}
}

func main() {
	lambda.Start(handler)
}

func jsonResponse(status int, body interface{}) (events.APIGatewayV2HTTPResponse, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error":"internal error"}`,
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(jsonBody),
	}, nil
}

func redirectReponse(status int, targetURL string) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Location": targetURL,
		},
	}, nil
}
