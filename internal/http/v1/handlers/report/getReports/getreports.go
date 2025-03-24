package getreports

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
	"google.golang.org/grpc/metadata"
)

type Response struct {
	resp.Response
	Reports []models.ReportResponse `json:"reports"`
}

type GetUsersClient interface {
	GetUsers(ctx context.Context, usersIds []int64) ([]models.User, error)
}

type GetReportsClient interface {
	GetReports(ctx context.Context, groupId, studentId, trainerId int64) ([]models.GrpcReport, []int64, error)
}

func New(reportCLient GetReportsClient, userClient GetUsersClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLoggerFromCtx(r.Context())

		cookies := r.CookiesNamed("authorization")
		if len(cookies) == 0 {
			log.Error(r.Context(), "cookies is empty")

			render.JSON(w, r, resp.Error("unauthorized"))

			return
		}

		var groupId int64
		var studentId int64
		var trainerId int64

		var err error
		if r.URL.Query().Has("group_id") {
			groupId, err = strconv.ParseInt(r.URL.Query().Get("group_id"), 10, 64)
			if err != nil {
				log.Error(r.Context(), "failed to parse user id", zap.Error(err))

				render.JSON(w, r, resp.Error(err.Error()))

				return
			}
		}

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

		accessToken := cookies[0].Value
		ctx := metadata.AppendToOutgoingContext(r.Context(), "authorization", accessToken)
		reports, ids, err := reportCLient.GetReports(ctx, groupId, studentId, trainerId)
		if err != nil {
			log.Error(r.Context(), "failed to get reports", zap.Error(err))

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		if len(reports) == 0 {
			render.JSON(w, r, Response{
				Response: resp.OK(),
				Reports:  []models.ReportResponse{},
			})
			return
		}

		users, err := userClient.GetUsers(ctx, ids)
		if err != nil {
			log.Error(r.Context(), "failed to get users", zap.Error(err))

			render.JSON(w, r, resp.Error(err.Error()))

			if errors.Is(err, clients.ErrUnauthenticated) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}

			return
		}

		usersMap := convert.ConvertToUserMap(users)
		reportsResp := make([]models.ReportResponse, len(reports))
		for i := range reports {
			reportsResp[i] = convert.ConvertReport(reports[i], usersMap)
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Reports:  reportsResp,
		})
	}
}
