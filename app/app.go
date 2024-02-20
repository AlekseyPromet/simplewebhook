package app

import (
	"AlekseyPromet/examples/simplewebhook/models"
	"AlekseyPromet/examples/simplewebhook/store"
	"context"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Service struct {
	port   string
	logger *zap.Logger
	store  store.Store
}

func NewService(cfg models.Config) (*Service, error) {
	var (
		err    error
		logger *zap.Logger
	)

	if cfg.Verbose && cfg.Debug {
		logger, err = zap.NewDevelopment(zap.AddStacktrace(zapcore.WarnLevel))
		if err != nil {
			return nil, fmt.Errorf("logger creation failed")
		}
	} else if cfg.Verbose {
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("logger creation failed")
		}
	} else {
		logger = zap.NewNop()
	}

	return &Service{
		port:   cfg.Port,
		logger: logger,
	}, fmt.Errorf("Service creation failed")
}

func (s *Service) Run(lc fx.Lifecycle) *http.Server {
	srv := &http.Server{Addr: "localhost:" + s.port}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			s.logger.Sugar().Infoln("Starting HTTP server at", srv.Addr)

			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.logger.Sugar().Infoln("server stopped", srv.Addr)
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
