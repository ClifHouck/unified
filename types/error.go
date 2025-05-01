package types

import (
	"time"
)

type Error struct {
	StatusCode  int       `json:"statusCode"`
	StatusName  string    `json:"statusName"`
	Message     string    `json:"message"`
	Timestamp   time.Time `json:"timestamp"`
	RequestPath string    `json:"requestPath"`
	RequestId   string    `json:"requestId"`
}
