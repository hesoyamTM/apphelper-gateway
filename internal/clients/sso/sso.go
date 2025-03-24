package sso

import (
	"context"

	"github.com/hesoyamTM/apphelper-gateway/internal/clients"
	"github.com/hesoyamTM/apphelper-gateway/internal/models"
	ssov1 "github.com/hesoyamTM/apphelper-protos/gen/go/sso"
	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log *logger.Logger
	api ssov1.AuthClient
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
		api: ssov1.NewAuthClient(cc),
	}, nil
}

func (c *Client) Register(ctx context.Context, name, surname, login, password string) (string, string, error) {
	res, err := c.api.Register(ctx, &ssov1.RegisterRequest{
		Name:     name,
		Surname:  surname,
		Login:    login,
		Password: password,
	})
	if err != nil {
		return "", "", err
	}

	return res.GetAccessToken(), res.GetRefreshToken(), nil
}

func (c *Client) Login(ctx context.Context, login, password string) (string, string, error) {
	res, err := c.api.Login(ctx, &ssov1.LoginRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return "", "", err
	}

	return res.GetAccessToken(), res.GetRefreshToken(), nil
}

func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	res, err := c.api.RefreshToken(ctx, &ssov1.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return "", "", err
	}

	return res.GetAccessToken(), res.GetRefreshToken(), nil
}

func (c *Client) GetUser(ctx context.Context, id int64) (models.User, error) {
	res, err := c.api.GetUser(ctx, &ssov1.GetUserRequest{
		UserId: id,
	})
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		Id:      id,
		Name:    res.GetName(),
		Surname: res.GetSurname(),
	}, nil
}

func (c *Client) GetUsers(ctx context.Context, ids []int64) ([]models.User, error) {
	res, err := c.api.GetUsers(ctx, &ssov1.GetUsersRequest{
		UserIds: ids,
	})
	if err != nil {
		return nil, err
	}

	users := make([]models.User, len(res.Users))
	for i := range res.Users {
		users[i].Id = res.Users[i].Id
		users[i].Name = res.Users[i].Name
		users[i].Surname = res.Users[i].Surname
	}

	return users, nil
}
