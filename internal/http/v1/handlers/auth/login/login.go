package login

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/hesoyamTM/apphelper-gateway/internal/lib/api/resp"
	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
	"go.uber.org/zap"
)

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	resp.Response
}

type LoginClient interface {
	Login(ctx context.Context, login, password string) (string, string, error)
}

func New(logClient LoginClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLoggerFromCtx(r.Context())

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error(r.Context(), "failed to decode request body", zap.Error(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		ctx := r.Context()
		accessToken, refreshToken, err := logClient.Login(ctx, req.Login, req.Password)
		if err != nil {
			log.Error(r.Context(), "failed to authorize user", zap.Error(err))

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		accessTokenCookie := http.Cookie{
			Name:     "authorization",
			Value:    fmt.Sprintf("Bearer %s", accessToken),
			HttpOnly: true,
			Domain:   "localhost",
		}
		refreshTokenCookie := http.Cookie{
			Name:     "refresh",
			Value:    refreshToken,
			HttpOnly: true,
			Domain:   "localhost",
		}

		http.SetCookie(w, &accessTokenCookie)
		http.SetCookie(w, &refreshTokenCookie)

		render.Status(r, http.StatusAccepted)
		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
