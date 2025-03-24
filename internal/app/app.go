package app

import (
	"context"
	"crypto/ecdsa"

	"github.com/hesoyamTM/apphelper-gateway/internal/app/httpapp"
	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
)

type App struct {
	log *logger.Logger

	HttpApp *httpapp.App
}

func New(
	ctx context.Context,
	hOpts httpapp.HttpOpts,
	clients httpapp.Clients,
	publicKey *ecdsa.PublicKey,
) *App {
	// TODO: services

	httpApp := httpapp.New(ctx, hOpts, clients, publicKey)

	return &App{
		HttpApp: httpApp,
	}
}
