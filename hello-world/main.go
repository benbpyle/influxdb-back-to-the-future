package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var client influxdb2.Client

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	writeClient := client.WriteAPIBlocking("Live Stream", "live-stream")

	locations := make([]string, 0)
	locations = append(locations,
		"Clocktower",
		"Diner",
		"Delorean",
		"Biff",
		"Tree",
		"Manure Truck")

	min := 0.0
	max := 2.0

	rand.Seed(time.Now().Unix())
	location := locations[rand.Intn(len(locations))]
	power := min + rand.Float64()*(max-min)

	p := influxdb2.NewPoint("strikes",
		map[string]string{"location": location},
		map[string]interface{}{"power": power},
		time.Now())

	writeClient.WritePoint(context.Background(), p)

	s := fmt.Sprintf("{ \"location\": \"%s\", \"power\": \"%s\" }", location, power)
	return events.APIGatewayProxyResponse{
		Body:       s,
		StatusCode: 201,
	}, nil
}

func main() {
	lambda.Start(handler)
}

func init() {
	client = influxdb2.NewClient("https://us-east-1-1.aws.cloud2.influxdata.com", "NYakOMbVSRZUHil-wmU9KJYdPNIJz7LKRY01pFIM-vN8rgIAcJWVfhMB-Hwv2KMMLV_kXZVnv3R647EwjlS6bA==")
}
