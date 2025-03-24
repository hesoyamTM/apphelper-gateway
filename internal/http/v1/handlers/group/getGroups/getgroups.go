package getgroups

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/hesoyamTM/apphelper-gateway/internal/clients"
	"github.com/hesoyamTM/apphelper-gateway/internal/lib/api/resp"
	"github.com/hesoyamTM/apphelper-gateway/internal/lib/convert"
	"github.com/hesoyamTM/apphelper-gateway/internal/models"
	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
	"go.uber.org/zap"
)

type Response struct {
	resp.Response
	Groups []models.GroupResponse `json:"groups"`
}

type GetUsersClient interface {
	GetUsers(ctx context.Context, usersIds []int64) ([]models.User, error)
}

type GetGroupsClient interface {
	GetGroups(ctx context.Context, studentId, trainerId int64) ([]models.GrpcGroup, []int64, error)
}

func New(groupsClient GetGroupsClient, userClient GetUsersClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLoggerFromCtx(r.Context())

		var studentId int64
		var trainerId int64

		var err error
		if r.URL.Query().Has("student_id") {
			studentId, err = strconv.ParseInt(r.URL.Query().Get("student_id"), 10, 64)
			if err != nil {
				log.Error(r.Context(), "failed to parse user id", zap.Error(err))

				render.JSON(w, r, resp.Error(err.Error()))

				return
			}
		}

		if r.URL.Query().Has("trainer_id") {
			trainerId, err = strconv.ParseInt(r.URL.Query().Get("trainer_id"), 10, 64)
			if err != nil {
				log.Error(r.Context(), "failed to parse user id", zap.Error(err))

				render.JSON(w, r, resp.Error(err.Error()))

				return
			}
		}

		groups, ids, err := groupsClient.GetGroups(r.Context(), studentId, trainerId)
		if err != nil {

			log.Error(r.Context(), "failed to get groups", zap.Error(err))

			if errors.Is(err, clients.ErrUnauthenticated) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		if len(groups) == 0 {
			render.JSON(w, r, Response{
				Response: resp.OK(),
				Groups:   []models.GroupResponse{},
			})
			return
		}

		users, err := userClient.GetUsers(r.Context(), ids)
		if err != nil {
			log.Error(r.Context(), "failed to get users", zap.Error(err))

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		usersMap := convert.ConvertToUserMap(users)
		groupResp := make([]models.GroupResponse, len(groups))
		for i := range groups {
			groupResp[i] = convert.ConvertGroup(groups[i], usersMap)
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Groups:   groupResp,
		})
	}
}
