package getuser

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/hesoyamTM/apphelper-gateway/internal/lib/api/resp"
	"github.com/hesoyamTM/apphelper-gateway/internal/models"
	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
	"go.uber.org/zap"
)

type Request struct {
	UserId int `json:"user_id"`
}

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

		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Error(r.Context(), "failed to parse id", zap.Error(err))

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		user, err := userClient.GetUser(r.Context(), int64(id))
		if err != nil {
			log.Error(r.Context(), "failed to get user", zap.Error(err))

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
			User:     user,
		})
	}
}
