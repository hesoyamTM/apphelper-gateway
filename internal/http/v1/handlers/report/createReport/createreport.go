package createreport

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/hesoyamTM/apphelper-gateway/internal/clients"
	"github.com/hesoyamTM/apphelper-gateway/internal/lib/api/resp"
	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
	"go.uber.org/zap"
)

type Request struct {
	GroupId     int64  `json:"group_id"`
	StudentId   int64  `json:"student_id"`
	TrainerId   int64  `json:"trainer_id"`
	Description string `json:"description"`
}

type Response struct {
	resp.Response
}

type CreteReportsClient interface {
	CreateReport(ctx context.Context, groupId, studentId, trainerId int64, desc string) error
}

func New(reportCLient CreteReportsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log := logger.GetLoggerFromCtx(r.Context())

		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error(r.Context(), "failed to decode request", zap.Error(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		if err := reportCLient.CreateReport(r.Context(), req.GroupId, req.StudentId, req.TrainerId, req.Description); err != nil {
			log.Error(r.Context(), "failed to create report", zap.Error(err))

			if errors.Is(err, clients.ErrUnauthenticated) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		render.JSON(w, r, resp.OK())
	}
}
