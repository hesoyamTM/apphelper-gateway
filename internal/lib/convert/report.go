package convert

import "github.com/hesoyamTM/apphelper-gateway/internal/models"

func ConvertReport(grpcRep models.GrpcReport, users map[int64]models.User) models.ReportResponse {
	return models.ReportResponse{
		GroupId:     grpcRep.GroupId,
		Student:     users[grpcRep.StudentId],
		Trainer:     users[grpcRep.TrainerId],
		Description: grpcRep.Description,
		Date:        grpcRep.Date,
	}
}
