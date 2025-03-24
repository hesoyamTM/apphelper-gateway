package refreshtoken

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/hesoyamTM/apphelper-gateway/internal/lib/api/resp"
	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
	"go.uber.org/zap"
)

type Response struct {
	resp.Response
}

type RefreshClient interface {
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
}

func New(refClient RefreshClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLoggerFromCtx(r.Context())

		cookies := r.CookiesNamed("refresh")
		if len(cookies) == 0 {
			log.Error(r.Context(), "failed to fetch cookies")

			http.Error(w, "unauthorized", http.StatusUnauthorized)
			render.JSON(w, r, resp.Error("unauthorized"))

			return
		}

		token := cookies[0].Value
		accessToken, refreshToken, err := refClient.RefreshToken(r.Context(), token)
		if err != nil {
			log.Error(r.Context(), "failed to refresh token", zap.Error(err))

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

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
