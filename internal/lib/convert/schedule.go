package convert

import "github.com/hesoyamTM/apphelper-gateway/internal/models"

func ConvertSchedule(grpcSched models.GrpcSchedule, users map[int64]models.User) models.ScheduleResponse {
	return models.ScheduleResponse{
		GroupName: grpcSched.GroupName,
		GroupId:   grpcSched.GroupId,
		Student:   users[grpcSched.StudentId],
		Trainer:   users[grpcSched.TrainerId],
		Date:      grpcSched.Date,
	}
}
