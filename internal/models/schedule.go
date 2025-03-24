package models

import (
	"time"
)

type GrpcSchedule struct {
	GroupName string
	GroupId   int64
	StudentId int64
	TrainerId int64
	Date      time.Time
}

type ScheduleResponse struct {
	GroupName string    `json:"group_name"`
	GroupId   int64     `json:"group_id"`
	Student   User      `json:"student"`
	Trainer   User      `json:"trainer"`
	Date      time.Time `json:"date"`
}
