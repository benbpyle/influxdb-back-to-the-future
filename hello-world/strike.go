package main

import "time"

type Strike struct {
	StrikeTime time.Time `json:"strikeTime"`
	Power      string    `json:"power"`
	Location   string    `json:"location"`
}
