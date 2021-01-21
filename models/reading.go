package models

import "time"

type Reading struct {
	Id        string      `json:"id"`
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	Alert     bool        `json:"alert"`
	Timestamp time.Time   `json:"timestamp"`
}

type ReadingSlice []Reading

func (r *ReadingSlice) Filter(f func(reading Reading) bool) *ReadingSlice {
	result := make(ReadingSlice, 0)
	for _, reading := range *r {
		if f(reading) {
			result = append(result, reading)
		}
	}
	return &result
}
