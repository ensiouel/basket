package app

import (
	"context"
	pb_static "github.com/ensiouel/basket-contract/gen/go/static/v1"
	"github.com/ensiouel/basket/internal/config"
	"github.com/ensiouel/basket/internal/service"
	"github.com/ensiouel/basket/internal/storage"
	"github.com/ensiouel/basket/internal/transport/rest"
	"github.com/ensiouel/basket/internal/transport/rest/handler"
	"github.com/ensiouel/basket/pkg/postgres"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	conf   config.Config
	logger *slog.Logger
}

func New() *App {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	logger := initLogger(conf.Logger)

	return &App{
		conf:   conf,
		logger: logger,
	}
}

func (app *App) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	pgConfig := postgres.Config{
		Host:     app.conf.Postgres.Host,
		Port:     app.conf.Postgres.Port,
		User:     app.conf.Postgres.User,
		Password: app.conf.Postgres.Password,
		DB:       app.conf.Postgres.DB,
	}
	pgClient, err := postgres.NewClient(ctx, pgConfig)
	if err != nil {
		app.logger.Error("cannot connect to postgres",
			slog.String("error", err.Error()),
		)
		return
	}

	grpcConn, err := grpc.Dial(app.conf.GRPC.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		app.logger.Error("cannot dial server",
			slog.String("error", err.Error()),
			slog.String("addr", app.conf.GRPC.Addr),
		)
	}
	defer grpcConn.Close()

	staticClient := pb_static.NewStaticClient(grpcConn)

	fileStorage := storage.NewFileStorage(pgClient)
	fileService := service.NewFileService(staticClient, fileStorage, app.conf.FileService.MaxFileSize)
	fileHandler := handler.NewFileHandler(fileService)

	server := rest.New(app.conf.Server, app.logger)

	app.logger.Info("starting http server", slog.String("addr", app.conf.Server.Addr))
	go func() {
		err = server.Handle(fileHandler).Run()
		if err != nil && err != http.ErrServerClosed {
			app.logger.Error("cannot run server",
				slog.String("error", err.Error()),
				slog.String("addr", app.conf.Server.Addr),
			)
		}
	}()

	<-ctx.Done()

	app.logger.Info("gracefully shutting down")
}

func initLogger(conf config.Logger) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: conf.Level,
	}))

	return logger
}
