package route

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Module("route",
	fx.Provide(
		fx.Annotate(NewMux,
			fx.ParamTags(`group:"routes"`),
			fx.As(fx.Self()),
			fx.As(new(http.Handler)),
		),
	),
	fx.Provide(newServer),
	fx.Decorate(decorateServer),
	fx.Invoke(func(*http.Server) {}),
)

func newServer(handler http.Handler) *http.Server {
	addr := net.JoinHostPort("", fmt.Sprint(viper.GetInt("port")))
	fmt.Printf("%T", handler)
	srv := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return &srv
}

func decorateServer(srv *http.Server, lc fx.Lifecycle) *http.Server {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				slog.Info("http server listening", "addr", srv.Addr)
				_ = srv.ListenAndServe()
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
