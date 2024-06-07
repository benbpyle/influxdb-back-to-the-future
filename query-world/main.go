package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/sirupsen/logrus"
)

var client influxdb2.Client

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	location := request.PathParameters["location"]
	query := fmt.Sprintf("from(bucket:\"live-stream\") |> range(start: -8h) |> filter(fn: (r) => r.location == \"%s\")", location)
	logrus.WithFields(logrus.Fields{
		"params": request.PathParameters,
		"query":  query,
	}).Debug("Just a debug request")

	queryAPI := client.QueryAPI("Live Stream")
	result, err := queryAPI.Query(
		context.Background(), query)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error with the query")

		return events.APIGatewayProxyResponse{
			Body:       "Error happened",
			StatusCode: 400,
		}, nil
	}

	powerContainer := &PowerContainer{}
	var readings []*Power

	for result.Next() {
		logrus.WithFields(logrus.Fields{
			"value": result.Record().Value(),
		}).Debug("Record value")

		power := newPower(result.Record())
		readings = append(readings, power)
	}

	powerContainer.Readings = readings
	b, _ := json.Marshal(powerContainer)

	return events.APIGatewayProxyResponse{
		Body:       string(b),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

func init() {
	client = influxdb2.NewClient("https://us-east-1-1.aws.cloud2.influxdata.com", "NYakOMbVSRZUHil-wmU9KJYdPNIJz7LKRY01pFIM-vN8rgIAcJWVfhMB-Hwv2KMMLV_kXZVnv3R647EwjlS6bA==")
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
}
