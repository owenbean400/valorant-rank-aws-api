package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"valorant-rank-api/domain/environment"
	"valorant-rank-api/domain/version"
	"valorant-rank-api/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Printf("req.Path: %s", req.RawPath)
	log.Printf("req.StageVariables: %+v", req.StageVariables)
	switch req.RawPath {
	case "/history":
		rank_data, err := service.GetValorantRankHistory(environment.GetPlayerPuuidEnv())

		if err != nil {
			return jsonResponse(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		} else {
			return jsonResponse(http.StatusOK, rank_data)
		}

	case "/update":
		err := service.UpdateDataWithAPI(environment.GetPlayerPuuidEnv())

		if err != nil {
			return jsonResponse(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		} else {
			return jsonResponse(http.StatusOK, map[string]string{"msg": "data updated"})
		}
	case "/version":
		return jsonResponse(http.StatusOK, map[string]string{"version": version.GetVersionNumber()})

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
