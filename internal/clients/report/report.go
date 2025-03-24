package report

import (
	"context"

	"github.com/hesoyamTM/apphelper-gateway/internal/clients"
	"github.com/hesoyamTM/apphelper-gateway/internal/models"
	reportv1 "github.com/hesoyamTM/apphelper-protos/gen/go/report"
	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	log *logger.Logger
	api reportv1.ReportClient
}

func New(ctx context.Context, addr string) (*Client, error) {
	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(clients.NewUIDInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		log: logger.GetLoggerFromCtx(ctx),
		api: reportv1.NewReportClient(cc),
	}, nil
}

func (c *Client) CreateReport(ctx context.Context, groupId, studentId, trainerId int64, desc string) error {
	_, err := c.api.CreateReport(ctx, &reportv1.CreateReportRequest{
		GroupId:     groupId,
		StudentId:   studentId,
		TrainerId:   trainerId,
		Description: desc,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return clients.ErrUnauthenticated
			}
		}
		return err
	}

	return nil
}

func (c *Client) GetReports(ctx context.Context, groupId, studentId, trainerId int64) ([]models.GrpcReport, []int64, error) {
	res, err := c.api.GetReports(ctx, &reportv1.GetReportsRequest{
		GroupId:   groupId,
		TrainerId: trainerId,
		StudentId: studentId,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, nil, clients.ErrUnauthenticated
			}
		}
		return nil, nil, err
	}

	reportsResp := res.GetReports()
	reports := make([]models.GrpcReport, len(reportsResp))
	userIds := make([]int64, 0, len(reportsResp))
	uniqueIds := make(map[int]struct{}, len(reportsResp))

	for i := range reportsResp {
		student := reportsResp[i].StudentId
		trainer := reportsResp[i].TrainerId
		group := reportsResp[i].GroupId

		reports[i] = models.GrpcReport{
			GroupId:     group,
			StudentId:   student,
			TrainerId:   trainer,
			Description: reportsResp[i].Description,
			Date:        reportsResp[i].Date.AsTime(),
		}

		if _, ok := uniqueIds[int(student)]; !ok {
			userIds = append(userIds, student)
		}
		if _, ok := uniqueIds[int(trainer)]; !ok {
			userIds = append(userIds, trainer)
		}
	}

	return reports, userIds, nil
}
