package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/pkg/v3/cobrautl"
)

var (
	version string
	mode    string
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	switch err := cmd.ExecuteContext(ctx); {
	case errors.As(err, new(*pflag.InvalidSyntaxError)):
	case errors.As(err, new(*pflag.InvalidValueError)):
		cobrautl.ExitWithError(cobrautl.ExitInvalidInput, err)
	}
}

var cmd = cobra.Command{
	Use:     filepath.Base(os.Args[0]),
	Version: version,
	Short:   "A Minecraft server management platform",
	PreRunE: setup,
	RunE:    run,
	PostRun: teardown,
}

func init() {
	cmd.Flags().StringP("config", "c", "/etc/herobrian/config.yml", "path to config file")
	cmd.Flags().StringP("environment", "e", "development", "environment name")
	cmd.Flags().IntP("log-level", "l", int(slog.LevelInfo), "default log level")
	cmd.Flags().IntP("port", "p", 8080, "http listener port")

	viper.SetDefault("log.level.stdlib", int(slog.LevelDebug))
	viper.SetDefault("log.level.http", int(slog.LevelWarn))

	viper.RegisterAlias("log.level.default", "log-level")
}

func setup(cmd *cobra.Command, _ []string) error {
	if config := cmd.Flags().Lookup("config"); config.Changed {
		viper.SetConfigFile(config.Value.String())
	} else {
		viper.SetConfigName("herobrian")
		viper.AddConfigPath("/etc")
		viper.AddConfigPath(".")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer("__", "."))
	viper.SetEnvPrefix("herobrian")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(viper.GetInt("log.level.default")),
	}))

	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(slog.Level(viper.GetInt("log.level.stdlib")))

	slog.Info("host started",
		"version", version,
		"mode", mode,
		"environment", viper.GetString("environment"),
	)

	return nil
}

func run(cmd *cobra.Command, _ []string) error {
	router := chi.NewRouter()
	router.Use(
		middleware.Recoverer,
		httplog.RequestLogger(slog.Default(), &httplog.Options{
			Level: slog.Level(viper.GetInt("log.level.http")),
		}),
	)

	secret := []byte(viper.GetString("jwt.secret"))
	auth := jwtauth.New(jwa.HS256.String(), secret, secret,
		jwt.WithAudience(viper.GetString("jwt.audience")),
		jwt.WithIssuer(viper.GetString("jwt.issuer")),
	)

	router.Use(jwtauth.Verifier(auth))

	var fileServer http.Handler
	if viper.GetString("environment") == "production" {
		fileServer = http.FileServer(http.Dir(viper.GetString("http.dir")))
	} else {
		url, _ := url.Parse("http://localhost:5173")
		fileServer = httputil.NewSingleHostReverseProxy(url)
	}

	router.Mount("/", fileServer)

	addr := net.JoinHostPort("", fmt.Sprint(viper.GetInt("port")))
	srv := http.Server{
		Addr:    addr,
		Handler: router,
	}

	go listen(&srv)

	<-cmd.Context().Done()

	slog.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func listen(srv *http.Server) {
	slog.Info("http server listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}

	slog.Debug("http server closed")
}

func teardown(*cobra.Command, []string) {
	slog.Info("goodbye :)")
}
