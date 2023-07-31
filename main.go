package main

import (
	"context"
	"net"
	"net/http"
	"pismo/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load()

	fx.New(
		fx.Provide(
			func() (*zap.Logger, error) {
				zapConfig := zap.NewProductionConfig()
				zapConfig.EncoderConfig.MessageKey = "message"

				return zapConfig.Build()
			},
			NewHTTPServer,
			fx.Annotate(
				NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),

			AsRoute(handlers.NewHealthcheckHandler),
		),
		fx.WithLogger(
			func(log *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: log}
			},
		),

		fx.Invoke(func(*http.Server) {}),
	).Run()
}

func NewHTTPServer(lc fx.Lifecycle, mux *mux.Router, log *zap.Logger) *http.Server {
	server := &http.Server{Addr: ":8080", Handler: mux}
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				ln, err := net.Listen("tcp", server.Addr)
				if err != nil {
					return err
				}

				log.Info("Starting HTTP server", zap.String("addr", server.Addr))
				go server.Serve(ln)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return server.Shutdown(ctx)
			},
		},
	)
	return server
}

func NewServeMux(handlers []handlers.Handler) *mux.Router {
	mux := mux.NewRouter()
	for _, h := range handlers {
		mux.Handle(h.Route(), h).Methods(h.Method()...)
	}

	return mux
}

func AsRoute(f interface{}) interface{} {
	return fx.Annotate(
		f,
		fx.As(new(handlers.Handler)),
		fx.ResultTags(`group:"routes"`),
	)
}