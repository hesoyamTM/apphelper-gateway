package addtogroup

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/hesoyamTM/apphelper-gateway/internal/clients"
	"github.com/hesoyamTM/apphelper-gateway/internal/lib/api/resp"
	"github.com/hesoyamTM/apphelper-sso/pkg/authorization"
	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
	"go.uber.org/zap"
)

type Request struct {
	Link string `json:"link"`
}

type AddToGroupClient interface {
	AddToGroup(ctx context.Context, studentId int64, link string) error
}

func New(groupClient AddToGroupClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLoggerFromCtx(r.Context())

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error(r.Context(), "failed to decode request", zap.Error(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		uid := r.Context().Value(authorization.Uid).(int64)
		err := groupClient.AddToGroup(r.Context(), uid, req.Link)
		if err != nil {
			log.Error(r.Context(), "failed to add to group", zap.Error(err))

			if errors.Is(err, clients.ErrUnauthenticated) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		render.JSON(w, r, resp.OK())
	}
}
