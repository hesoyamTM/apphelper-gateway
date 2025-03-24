package getusers

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

type Response struct {
	resp.Response
	Users []models.User `json:"users"`
}

type GetUserClient interface {
	GetUsers(ctx context.Context, ids []int64) ([]models.User, error)
}

func New(userClient GetUserClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLoggerFromCtx(r.Context())

		r.ParseForm()
		ids := r.Form["ids"]

		userIds := make([]int64, len(ids))
		for i, id := range ids {
			var err error
			userIds[i], err = strconv.ParseInt(id, 10, 64)
			if err != nil {
				log.Error(r.Context(), "failed to parse user id", zap.Error(err))

				render.JSON(w, r, resp.Error(err.Error()))

				return
			}
		}

		users, err := userClient.GetUsers(r.Context(), userIds)
		if err != nil {
			log.Error(r.Context(), "failed to get user", zap.Error(err))

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Users:    users,
		})
	}
}
