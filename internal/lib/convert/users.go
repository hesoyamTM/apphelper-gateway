package convert

import "github.com/hesoyamTM/apphelper-gateway/internal/models"

func ConvertToUserMap(users []models.User) map[int64]models.User {
	usersMap := make(map[int64]models.User, len(users))

	for _, user := range users {
		usersMap[user.Id] = user
	}

	return usersMap
}
