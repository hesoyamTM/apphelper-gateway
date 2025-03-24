package httpapp

import (
	"context"
	"crypto/ecdsa"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hesoyamTM/apphelper-gateway/internal/clients/report"
	"github.com/hesoyamTM/apphelper-gateway/internal/clients/schedule"
	"github.com/hesoyamTM/apphelper-gateway/internal/clients/sso"
	"github.com/hesoyamTM/apphelper-gateway/internal/http/v1/handlers/auth/getsession"
	"github.com/hesoyamTM/apphelper-gateway/internal/http/v1/handlers/auth/login"
	refreshtoken "github.com/hesoyamTM/apphelper-gateway/internal/http/v1/handlers/auth/refreshToken"
	"github.com/hesoyamTM/apphelper-gateway/internal/http/v1/handlers/auth/register"
	addtogroup "github.com/hesoyamTM/apphelper-gateway/internal/http/v1/handlers/group/addToGroup"
	creategroup "github.com/hesoyamTM/apphelper-gateway/internal/http/v1/handlers/group/createGroup"
	getgroups "github.com/hesoyamTM/apphelper-gateway/internal/http/v1/handlers/group/getGroups"
	createreport "github.com/hesoyamTM/apphelper-gateway/internal/http/v1/handlers/report/createReport"
	getreports "github.com/hesoyamTM/apphelper-gateway/internal/http/v1/handlers/report/getReports"
	createschedule "github.com/hesoyamTM/apphelper-gateway/internal/http/v1/handlers/schedule/createSchedule"
	createscheduleforgroup "github.com/hesoyamTM/apphelper-gateway/internal/http/v1/handlers/schedule/createScheduleForGroup"
	getschedules "github.com/hesoyamTM/apphelper-gateway/internal/http/v1/handlers/schedule/getSchedules"
	getuser "github.com/hesoyamTM/apphelper-gateway/internal/http/v1/handlers/user/getUser"
	getusers "github.com/hesoyamTM/apphelper-gateway/internal/http/v1/handlers/user/getUsers"
	"github.com/hesoyamTM/apphelper-gateway/internal/http/v1/middlewares"
	"github.com/hesoyamTM/apphelper-sso/pkg/authorization"
	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
	"go.uber.org/zap"
)

var (
	authMethods = map[string]bool{
		"/getSession": true,

		"/createReport": true,
		"/getReports":   true,

		"/createScheduleForGroup": true,
		"/createSchedule":         true,
		"/getSchedules":           true,

		"/createGroup": true,
		"/getGroups":   true,
		"/addToGroup":  true,
	}
)

type App struct {
	log        *logger.Logger
	router     *chi.Mux
	httpServer *http.Server
}

type HttpOpts struct {
	Addr        string
	Timeout     time.Duration
	TimeoutIdle time.Duration
}

type Clients struct {
	AuthAddrs    string
	ReportAddr   string
	ScheduleAddr string
}

func New(ctx context.Context, hOtps HttpOpts, clients Clients, privKey *ecdsa.PublicKey) *App {
	log := logger.GetLoggerFromCtx(ctx)
	router := chi.NewRouter()

	router.Use(middlewares.Cors)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middlewares.RequestLogging)
	router.Use(logger.LoggingMiddleware(ctx))
	router.Use(authorization.NewAuthMiddleware(authMethods, privKey))

	authClient, err := sso.New(ctx, clients.AuthAddrs)
	if err != nil {
		log.Error(ctx, "failed to connect to sso client", zap.Error(err))
	}
	repClient, err := report.New(ctx, clients.ReportAddr)
	if err != nil {
		log.Error(ctx, "failed to connect to report client", zap.Error(err))
	}
	schedClient, err := schedule.New(ctx, clients.ScheduleAddr)
	if err != nil {
		log.Error(ctx, "failed to connect to report client", zap.Error(err))
	}

	//sso
	router.Post("/register", register.New(authClient))
	router.Post("/login", login.New(authClient))
	router.Head("/refreshToken", refreshtoken.New(authClient))
	router.Get("/getSession", getsession.New(authClient))
	//users
	router.Get("/getUser", getuser.New(authClient))
	router.Get("/getUsers", getusers.New(authClient))
	//reports
	router.Post("/createReport", createreport.New(repClient))
	router.Get("/getReports", getreports.New(repClient, authClient))
	//schedules
	router.Post("/createSchedule", createschedule.New(schedClient))
	router.Post("/createScheduleForGroup", createscheduleforgroup.New(schedClient))
	router.Get("/getSchedules", getschedules.New(schedClient, authClient))
	//groups
	router.Post("/createGroup", creategroup.New(schedClient))
	router.Get("/getGroups", getgroups.New(schedClient, authClient))
	router.Patch("/addToGroup", addtogroup.New(schedClient))

	httpServer := &http.Server{
		Addr:         hOtps.Addr,
		Handler:      router,
		ReadTimeout:  hOtps.Timeout,
		WriteTimeout: hOtps.Timeout,
		IdleTimeout:  hOtps.TimeoutIdle,
	}

	return &App{
		log:        log,
		router:     router,
		httpServer: httpServer,
	}
}

func (a *App) MustRun(ctx context.Context) {
	a.log.Info(ctx, "http server is running")

	if err := a.httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (a *App) Stop(ctx context.Context) {
	a.log.Info(ctx, "server is stopped")

	if err := a.httpServer.Shutdown(context.Background()); err != nil {
		a.log.Error(ctx, "", zap.Error(err))
	}
}
