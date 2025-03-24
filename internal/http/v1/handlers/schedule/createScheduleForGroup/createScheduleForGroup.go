package createscheduleforgroup

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/hesoyamTM/apphelper-gateway/internal/clients"
	"github.com/hesoyamTM/apphelper-gateway/internal/lib/api/resp"
	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
	"go.uber.org/zap"
)

type Request struct {
	GroupId   int64     `json:"group_id"`
	TrainerId int64     `json:"trainer_id"`
	Date      time.Time `json:"date"`
}

type CreateScheduleClient interface {
	CreateScheduleForGroup(ctx context.Context, groupId, trainerId int64, date time.Time) error
}

func New(scheduleClient CreateScheduleClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLoggerFromCtx(r.Context())

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error(r.Context(), "failed to decode request", zap.Error(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		err := scheduleClient.CreateScheduleForGroup(r.Context(), req.GroupId, req.TrainerId, req.Date)
		if err != nil {
			log.Error(r.Context(), "failed to create schedule", zap.Error(err))

			if errors.Is(err, clients.ErrUnauthenticated) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		render.JSON(w, r, resp.OK())
	}
}
