package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/bdreece/herobrian/internal/database"
	"github.com/bdreece/herobrian/internal/identity"
	"github.com/bdreece/herobrian/pkg/minecraft"
	"github.com/go-crypt/crypt"
	"github.com/go-crypt/crypt/algorithm/argon2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/pkg/v3/cobrautl"
	echovalidator "gopkg.in/bdreece/echo-validator.v1"
	_ "modernc.org/sqlite"
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

	viper.SetOptions(
		viper.KeyDelimiter(":"),
		viper.EnvKeyReplacer(strings.NewReplacer("__", ":")),
	)

	viper.SetDefault("app.root_dir", "/usr/share/herobrian/www")
	viper.SetDefault("app.proxy_url", "http://localhost:5173")
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

	viper.SetEnvPrefix("herobrian")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil && !errors.As(err, new(viper.ConfigFileNotFoundError)) {
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
	db, err := sql.Open("sqlite", viper.GetString("sqlite.dsn"))
	if err != nil {
		return err
	}

	defer db.Close()

	queries := database.New(db)
	if _, err := queries.ApplySchema(cmd.Context()); err != nil {
		return err
	}

	slog.Debug("opened sqlite connection")

	awsConfig, err := config.LoadDefaultConfig(cmd.Context())
	if err != nil {
		return err
	}

	slog.Debug("loaded AWS config")

	ec2 := ec2.NewFromConfig(awsConfig)
	slog.Debug("configured EC2 client")

	accessTokenSigner, err := identity.NewTokenSigner(identity.AccessToken)
	if err != nil {
		return err
	}

	inviteTokenSigner, err := identity.NewTokenSigner(identity.InviteToken)
	if err != nil {
		return err
	}

	slog.Debug("created token signers")

	decoder := crypt.NewDecoder()
	if err = argon2.RegisterDecoderArgon2id(decoder); err != nil {
		return err
	}

	hasher, err := argon2.New(
		argon2.WithProfileRFC9106LowMemory(),
	)
	if err != nil {
		return err
	}

	slog.Debug("created password hasher + decoder")

	hostConfigs := map[string]minecraft.HostConfig{}
	if err := viper.UnmarshalKey("hosts", &hostConfigs); err != nil {
		return err
	}

	provider := minecraft.EC2Provider{
		Client: ec2,
	}

	e := echo.New()
	e.Validator = echovalidator.Default

	e.POST("/login", identity.NewLoginHandler(queries, accessTokenSigner, decoder))
	e.POST("/identity/activate", identity.NewActivateHandler())

	hosts := e.Group("/api/host")
	hosts.GET("/", minecraft.NewHostHandler(&provider, hostConfigs))

	e.Use(
		middleware.Recover(),
	)

	env := viper.GetString("environment")
	var appServer echo.MiddlewareFunc
	if env == "production" {
		appServer = middleware.StaticWithConfig(middleware.StaticConfig{
			Root:  viper.GetString("app.root_dir"),
			HTML5: true,
		})
	} else {
		url, _ := url.Parse(viper.GetString("app.proxy_url"))
		balancer := middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
			{URL: url},
		})

		appServer = middleware.ProxyWithConfig(middleware.ProxyConfig{
			Balancer: balancer,
		})
	}

	e.Use(appServer)

	addr := net.JoinHostPort("", fmt.Sprint(viper.GetInt("port")))
	srv := http.Server{
		Addr: addr,
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
