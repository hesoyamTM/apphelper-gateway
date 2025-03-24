package models

import (
	"time"
)

type GrpcReport struct {
	GroupId     int64
	StudentId   int64
	TrainerId   int64
	Description string
	Date        time.Time
}

type ReportResponse struct {
	GroupId     int64     `json:"group_id"`
	Student     User      `json:"student"`
	Trainer     User      `json:"trainer"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}
