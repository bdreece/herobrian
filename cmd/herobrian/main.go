package main

import (
	"context"
	"crypto/rand"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/r3labs/sse/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/pkg/v3/cobrautl"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	_ "modernc.org/sqlite"

	"github.com/bdreece/herobrian/internal/database"
	"github.com/bdreece/herobrian/internal/route"
	"github.com/bdreece/herobrian/internal/security"
	"github.com/bdreece/herobrian/pkg/minecraft"
	"github.com/bdreece/herobrian/pkg/user"
)

var (
	version string
	mode    string

	app *fx.App

	cmd = cobra.Command{
		Use:      filepath.Base(os.Args[0]),
		Version:  version,
		Short:    "A Minecraft server management platform",
		PreRunE:  setup,
		RunE:     run,
		PostRunE: teardown,
	}
)

func init() {
	cmd.Flags().StringP("config", "c", "/etc/herobrian/config.yml", "path to config file")
	cmd.Flags().StringP("environment", "e", "development", "environment name")
	cmd.Flags().IntP("log-level", "l", int(slog.LevelInfo), "default log level")
	cmd.Flags().IntP("port", "p", 8080, "http listener port")

	viper.SetOptions(
		viper.KeyDelimiter(":"),
		viper.EnvKeyReplacer(strings.NewReplacer("__", ":")),
	)

	viper.SetDefault("aws:region", "us-east-2")

	viper.SetDefault("http:cookie:refresh", "herobrian-rt")
	viper.SetDefault("http:root_dir", "/usr/share/herobrian/www")
	viper.SetDefault("http:proxy_url", "http://localhost:5173")

	viper.SetDefault("jwt:default:aud", "herobrian.bdreece.dev")
	viper.SetDefault("jwt:default:iss", "herobrian.bdreece.dev")
	viper.SetDefault("jwt:access:lifetime", time.Hour)
	viper.SetDefault("jwt:access:secret", rand.Text())
	viper.SetDefault("jwt:invite:lifetime", time.Hour)
	viper.SetDefault("jwt:invite:secret", rand.Text())

	viper.SetDefault("log:level:fx", int(slog.LevelDebug))
	viper.SetDefault("log:level:stdlib", int(slog.LevelDebug))
	viper.SetDefault("log:level:http", int(slog.LevelWarn))

	viper.SetDefault("sqlite:dsn", "file:///usr/lib/herobrian/db.sqlite3")
	viper.SetDefault("sqlite:admin:first_name", "Steve")
	viper.SetDefault("sqlite:admin:last_name", "Minecraft")
	viper.SetDefault("sqlite:admin:display_name", "admin")
	viper.SetDefault("sqlite:admin:password", "password")
	viper.SetDefault("sqlite:admin:picture_url", "https://i.imgflip.com/2/k2klk.jpg")
	viper.SetDefault("sqlite:admin:role", "swashbuckler")

	viper.RegisterAlias("log:level:default", "log-level")
}

func setup(cmd *cobra.Command, _ []string) error {
	if config := cmd.Flags().Lookup("config"); config.Changed {
		viper.SetConfigFile(config.Value.String())
	} else {
		viper.SetConfigName("herobrian")
		viper.AddConfigPath("/etc")
		viper.AddConfigPath(".")
	}

	viper.SetEnvPrefix("herobrian")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil && !errors.As(err, new(viper.ConfigFileNotFoundError)) {
		return err
	}

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	environment := viper.GetString("environment")
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(viper.GetInt("log:level:default")),
	}))

	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(slog.Level(viper.GetInt("log:level:stdlib")))
	slog.Info("host started",
		"environment", environment,
		"version", version,
		"mode", mode,
	)

	return nil
}

func run(cmd *cobra.Command, _ []string) error {
	var logOption fx.Option

	if mode == "debug" {
		logOption = fx.WithLogger(newFxLogger)
	} else {
		logOption = fx.NopLogger
	}

	app = fx.New(
		logOption,
		fx.Supply(slog.Default()),
		fx.Provide(newAwsConfig, ec2.NewFromConfig),
		fx.Provide(func() *sse.Server {
			srv := sse.New()

			srv.AutoStream = true
			srv.EventTTL = time.Hour

			return srv
		}),
		database.Module,
		security.Module,
		user.Module,
		minecraft.Module,
		route.Module,
	)

	if err := app.Start(cmd.Context()); err != nil {
		return err
	}

	<-cmd.Context().Done()
	return nil
}

func teardown(cmd *cobra.Command, _ []string) error {
	return app.Stop(context.Background())
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	switch err := cmd.ExecuteContext(ctx); {
	case errors.As(err, new(*pflag.InvalidSyntaxError)):
	case errors.As(err, new(*pflag.InvalidValueError)):
		cobrautl.ExitWithError(cobrautl.ExitInvalidInput, err)
	}
}

func newAwsConfig() (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(viper.GetString("aws:region")),
	)
}

func newFxLogger() fxevent.Logger {
	logger := fxevent.SlogLogger{
		Logger: slog.Default().With("scope", "fx"),
	}

	logger.UseLogLevel(slog.Level(viper.GetInt("log:level:fx")))

	return &logger
}
