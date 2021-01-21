package models

import "time"

type Reading struct {
	Id        string      `json:"id"`
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	Alert     bool        `json:"alert"`
	Timestamp time.Time   `json:"timestamp"`
}
