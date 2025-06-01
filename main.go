package main

import (
	"context"
	"encoding/json"
	"net/http"
	"valorant-rank-api/domain/environment"
	"valorant-rank-api/domain/version"
	"valorant-rank-api/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.Path {
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

func jsonResponse(status int, body interface{}) (events.APIGatewayProxyResponse, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error":"internal error"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(jsonBody),
	}, nil
}
