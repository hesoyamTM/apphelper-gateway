package convert

import "github.com/hesoyamTM/apphelper-gateway/internal/models"

func ConvertGroup(grpcGroup models.GrpcGroup, users map[int64]models.User) models.GroupResponse {
	studentIds := grpcGroup.StudentIds
	students := make([]models.User, len(studentIds))

	for i, id := range studentIds {
		students[i] = users[id]
	}

	return models.GroupResponse{
		Id:       grpcGroup.Id,
		Students: students,
		Trainer:  users[grpcGroup.TrainerId],
		Name:     grpcGroup.Name,
		Link:     grpcGroup.Link,
	}
}
