package schedule

import (
	"context"
	"time"

	"github.com/hesoyamTM/apphelper-gateway/internal/clients"
	"github.com/hesoyamTM/apphelper-gateway/internal/models"
	schedulev1 "github.com/hesoyamTM/apphelper-protos/gen/go/schedule"
	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	log *logger.Logger
	api schedulev1.ScheduleClient
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
		api: schedulev1.NewScheduleClient(cc),
	}, nil
}

func (c *Client) CreateSchedule(ctx context.Context, groupId, studentId, trainerId int64, date time.Time) error {
	_, err := c.api.CreateSchedule(ctx, &schedulev1.CreateScheduleRequest{
		GroupId:   groupId,
		StudentId: studentId,
		TrainerId: trainerId,
		Date:      timestamppb.New(date),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) CreateScheduleForGroup(ctx context.Context, groupId, trainerId int64, date time.Time) error {
	_, err := c.api.CreateScheduleForGroup(ctx, &schedulev1.CreateScheduleForGroupRequest{
		GroupId:   groupId,
		TrainerId: trainerId,
		Date:      timestamppb.New(date),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetSchedules(ctx context.Context, studentId, trainerId int64) ([]models.GrpcSchedule, []int64, error) {
	resp, err := c.api.GetSchedule(ctx, &schedulev1.GetSchedulesRequest{
		TrainerId: trainerId,
		StudentId: studentId,
	})
	if err != nil {
		return nil, nil, err
	}

	schedules := resp.GetSchedules()
	schedulesResp := make([]models.GrpcSchedule, len(schedules))
	ids := make([]int64, 0, len(schedules))
	uniqueIds := make(map[int64]struct{}, len(schedules))

	for i := range schedules {
		student := schedules[i].StudentId
		trainer := schedules[i].TrainerId

		schedulesResp[i].GroupName = schedules[i].GroupName
		schedulesResp[i].GroupId = schedules[i].GroupId
		schedulesResp[i].StudentId = schedules[i].StudentId
		schedulesResp[i].TrainerId = schedules[i].TrainerId
		schedulesResp[i].Date = schedules[i].Date.AsTime()

		if _, ok := uniqueIds[student]; !ok {
			ids = append(ids, student)
		}
		if _, ok := uniqueIds[trainer]; !ok {
			ids = append(ids, trainer)
		}
	}

	return schedulesResp, ids, nil
}

func (c *Client) CreateGroup(ctx context.Context, trainerId int64, name string) error {
	_, err := c.api.CreateGroup(ctx, &schedulev1.CreateGroupRequest{
		TrainerId: trainerId,
		Name:      name,
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

func (c *Client) GetGroups(ctx context.Context, studentId, trainerId int64) ([]models.GrpcGroup, []int64, error) {
	resp, err := c.api.GetGroups(ctx, &schedulev1.GetGroupsRequest{
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

	groups := resp.GetGroups()
	groupsResp := make([]models.GrpcGroup, len(groups))
	ids := make([]int64, 0, len(groups))
	uniqueIds := make(map[int64]struct{}, len(groups))

	for i := range groups {
		students := groups[i].Students
		trainer := groups[i].TrainerId

		groupsResp[i].Id = groups[i].GetId()
		groupsResp[i].Name = groups[i].GetName()
		groupsResp[i].StudentIds = groups[i].GetStudents()
		groupsResp[i].TrainerId = groups[i].GetTrainerId()
		groupsResp[i].Link = groups[i].GetLink()

		for _, student := range students {
			if _, ok := uniqueIds[student]; !ok {
				ids = append(ids, student)
			}
		}
		if _, ok := uniqueIds[trainer]; !ok {
			ids = append(ids, trainer)
		}
	}

	return groupsResp, ids, nil
}

func (c *Client) AddToGroup(ctx context.Context, studentId int64, link string) error {
	_, err := c.api.AddToGroup(ctx, &schedulev1.AddToGroupRequest{
		StudentId: studentId,
		Link:      link,
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
