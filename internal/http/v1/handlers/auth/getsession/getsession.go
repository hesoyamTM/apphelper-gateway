package getsession

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/hesoyamTM/apphelper-gateway/internal/lib/api/resp"
	"github.com/hesoyamTM/apphelper-gateway/internal/models"
	"github.com/hesoyamTM/apphelper-sso/pkg/authorization"
	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
	"go.uber.org/zap"
)

type Response struct {
	resp.Response
	User models.User `json:"user"`
}

type GetUserClient interface {
	GetUser(ctx context.Context, id int64) (models.User, error)
}

func New(userClient GetUserClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLoggerFromCtx(r.Context())

		uid := r.Context().Value(authorization.Uid).(int64)
		user, err := userClient.GetUser(r.Context(), uid)
		if err != nil {
			log.Error(r.Context(), "failed to authorize user", zap.Error(err))

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
			User:     user,
		})
	}
}
