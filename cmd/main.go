package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"os/signal"
	"syscall"

	"github.com/hesoyamTM/apphelper-gateway/internal/app"
	"github.com/hesoyamTM/apphelper-gateway/internal/app/httpapp"
	"github.com/hesoyamTM/apphelper-gateway/internal/config"
	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
)

const (
	localEnv = "local"
	devEnv   = "dev"
	prodEnv  = "prod"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()

	ctx, err := logger.New(ctx, cfg.Env)
	if err != nil {
		panic(err)
	}

	log := logger.GetLoggerFromCtx(ctx)
	log.Debug(ctx, "logger is working")

	hOpts := httpapp.HttpOpts{
		Addr:        cfg.Http.Address,
		Timeout:     cfg.Http.Timeout,
		TimeoutIdle: cfg.Http.IdleTimeout,
	}

	clients := httpapp.Clients{
		AuthAddrs:    cfg.GrpcClients.SsoAddress,
		ReportAddr:   cfg.GrpcClients.ReportAddress,
		ScheduleAddr: cfg.GrpcClients.ScheduleAddress,
	}

	publicKey := decodePubKey(cfg.PublicKey)

	application := app.New(ctx, hOpts, clients, publicKey)
	go application.HttpApp.MustRun(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	application.HttpApp.Stop(ctx)
	log.Info(ctx, "application stopped")
}

func decodePubKey(pemEncodedPub string) *ecdsa.PublicKey {

	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return publicKey
}
