package creategroup

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
	Name      string `json:"name"`
	TrainerId int64  `json:"trainer_id"`
}

type CreateGroupClient interface {
	CreateGroup(ctx context.Context, trainerId int64, name string) error
}

func New(groupClient CreateGroupClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLoggerFromCtx(r.Context())

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error(r.Context(), "failed to decode request", zap.Error(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		err := groupClient.CreateGroup(r.Context(), req.TrainerId, req.Name)
		if err != nil {
			log.Error(r.Context(), "failed to create group", zap.Error(err))

			if errors.Is(err, clients.ErrUnauthenticated) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		render.JSON(w, r, resp.OK())
	}
}
