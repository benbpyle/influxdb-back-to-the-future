package main

import (
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/query"
)

type PowerContainer struct {
	Readings []*Power `json:"readings"`
}

type Power struct {
	Time  time.Time `json:"readingTime"`
	Power float64   `json:"readingPower"`
}

func newPower(record *query.FluxRecord) *Power {
	value, ok := record.Value().(float64)

	if !ok {
		return &Power{
			Time:  record.Time(),
			Power: 0.0,
		}
	}

	return &Power{
		Time:  record.Time(),
		Power: value,
	}
}
