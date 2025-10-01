package models

import "time"

type Count struct {
	Count int64 `json:"count,omitempty"`
}

type CountTime struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int64     `json:"count"`
}
