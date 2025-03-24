package getschedules

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
	Schedules []models.ScheduleResponse `json:"schedules"`
}

type GetUsersClient interface {
	GetUsers(ctx context.Context, usersIds []int64) ([]models.User, error)
}

type GetSchedulesClient interface {
	GetSchedules(ctx context.Context, studentId, trainerId int64) ([]models.GrpcSchedule, []int64, error)
}

func New(scheduleClient GetSchedulesClient, userClient GetUsersClient) http.HandlerFunc {
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

		schedules, ids, err := scheduleClient.GetSchedules(r.Context(), studentId, trainerId)
		if err != nil {
			log.Error(r.Context(), "failed to get schedules", zap.Error(err))

			if errors.Is(err, clients.ErrUnauthenticated) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		if len(schedules) == 0 {
			render.JSON(w, r, Response{
				Response:  resp.OK(),
				Schedules: []models.ScheduleResponse{},
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
		schedulesResp := make([]models.ScheduleResponse, len(schedules))
		for i := range schedules {
			schedulesResp[i] = convert.ConvertSchedule(schedules[i], usersMap)
		}

		render.JSON(w, r, Response{
			Response:  resp.OK(),
			Schedules: schedulesResp,
		})
	}
}
