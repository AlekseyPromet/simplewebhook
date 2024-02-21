package app

import (
	"AlekseyPromet/examples/simplewebhook/models"
	"AlekseyPromet/examples/simplewebhook/store"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Service struct {
	port   string
	logger *zap.Logger
	store  store.Store
}

func NewService(cfg *models.Config) (*Service, error) {
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
			logger: logger.Named("WHOOK"),
			store:  store.NewTestStore(*logger),
		},
		nil
}

func (s *Service) GetFxLogger() fxevent.Logger {
	return &fxevent.ZapLogger{Logger: s.logger}
}

func (s *Service) GetServeMux() *http.ServeMux {

	mux := http.NewServeMux()

	middleware := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
	}

	middlewareError := func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	mux.HandleFunc("POST /invoke", func(w http.ResponseWriter, r *http.Request) {
		middleware(w, r)

		source := models.Source{}

		if err := json.NewDecoder(r.Body).Decode(&source); err != nil {
			middlewareError(w, r, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancel()

		key, err := s.store.Create(ctx, source)
		if err != nil {
			middlewareError(w, r, err)
			return
		}

		if err := json.NewEncoder(w).Encode(key); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})

	mux.HandleFunc("POST /get", func(w http.ResponseWriter, r *http.Request) {
		middleware(w, r)

		source := models.ApiKey{}

		if err := json.NewDecoder(r.Body).Decode(&source); err != nil {
			middlewareError(w, r, err)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancel()

		response, err := s.store.Get(ctx, source.Key)
		if err != nil {
			middlewareError(w, r, err)
			return
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})

	mux.HandleFunc("POST /delete", func(w http.ResponseWriter, r *http.Request) {
		middleware(w, r)

		source := models.ApiKey{}

		if err := json.NewDecoder(r.Body).Decode(&source); err != nil {
			middlewareError(w, r, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancel()

		response, err := s.store.Delete(ctx, source.Key)
		if err != nil {
			middlewareError(w, r, err)
			return
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})

	return mux
}

func (s *Service) Run(lc fx.Lifecycle) *http.Server {

	srv := &http.Server{
		Addr:    "localhost:" + s.port,
		Handler: s.GetServeMux(),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			s.logger.Sugar().Infoln("Starting HTTP server at", srv.Addr)

			go srv.Serve(ln)

			go s.WebhookCycle(context.Background())

			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.logger.Sugar().Infoln("server stopped", srv.Addr)
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
