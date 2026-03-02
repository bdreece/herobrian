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
	"github.com/bdreece/herobrian/internal/database"
	"github.com/bdreece/herobrian/internal/identity"
	"github.com/bdreece/herobrian/internal/route"
	"github.com/bdreece/herobrian/pkg/minecraft"
	"github.com/bdreece/herobrian/pkg/user"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/pkg/v3/cobrautl"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	_ "modernc.org/sqlite"
)

var (
	version string
	mode    string
)

var cmd = cobra.Command{
	Use:     filepath.Base(os.Args[0]),
	Version: version,
	Short:   "A Minecraft server management platform",
	PreRunE: setup,
	Run:     run,
}

func init() {
	cmd.Flags().StringP("config", "c", "/etc/herobrian/config.yml", "path to config file")
	cmd.Flags().StringP("environment", "e", "development", "environment name")
	cmd.Flags().IntP("log-level", "l", int(slog.LevelInfo), "default log level")
	cmd.Flags().IntP("port", "p", 8080, "http listener port")

	viper.SetOptions(
		viper.KeyDelimiter(":"),
		viper.EnvKeyReplacer(strings.NewReplacer("__", ":")),
	)

	viper.SetDefault("app:root_dir", "/usr/share/herobrian/www")
	viper.SetDefault("app:proxy_url", "http://localhost:5173")

	viper.SetDefault("aws:region", "us-east-2")

	viper.SetDefault("jwt:default:aud", "herobrian.bdreece.dev")
	viper.SetDefault("jwt:default:iss", "herobrian.bdreece.dev")
	viper.SetDefault("jwt:access:lifetime", time.Hour)
	viper.SetDefault("jwt:access:secret", rand.Text())
	viper.SetDefault("jwt:invite:lifetime", time.Hour)
	viper.SetDefault("jwt:invite:secret", rand.Text())

	viper.SetDefault("log:level:fx", int(slog.LevelInfo))
	viper.SetDefault("log:level:stdlib", int(slog.LevelDebug))
	viper.SetDefault("log:level:http", int(slog.LevelWarn))

	viper.SetDefault("sqlite:dsn", "file:///usr/lib/herobrian/db.sqlite3")
	viper.SetDefault("sqlite:admin:first_name", "Steve")
	viper.SetDefault("sqlite:admin:last_name", "Minecraft")
	viper.SetDefault("sqlite:admin:display_name", "admin")
	viper.SetDefault("sqlite:admin:password", "password")
	viper.SetDefault("sqlite:admin:picture_url", "https://i.imgflip.com/2/k2klk.jpg")

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

func run(cmd *cobra.Command, _ []string) {

	var logOption fx.Option
	if viper.GetInt("log:level:fx") >= 0 {
		logOption = fx.NopLogger
	} else {
		logOption = fx.WithLogger(func() fxevent.Logger {
			return &fxevent.SlogLogger{Logger: slog.Default()}
		})
	}

	awsOptions := fx.Options(
		fx.Provide(func() (aws.Config, error) {
			return config.LoadDefaultConfig(context.TODO(),
				config.WithRegion(viper.GetString("aws:region")),
			)
		}),
		fx.Provide(ec2.NewFromConfig),
	)

	fx.New(
		logOption,
		awsOptions,
		database.Module,
		identity.Module,
		user.Module,
		minecraft.Module,
		route.Module,
	).Run()
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
