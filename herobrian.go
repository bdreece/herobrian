package herobrian

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"maps"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	echorenderer "github.com/bdreece/echo-renderer"
	echovalidator "github.com/bdreece/echo-validator"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"golang.org/x/crypto/bcrypt"

	"github.com/bdreece/herobrian/internal/controller"
	"github.com/bdreece/herobrian/internal/logger"
	"github.com/bdreece/herobrian/internal/middleware"
	"github.com/bdreece/herobrian/internal/router"
	"github.com/bdreece/herobrian/pkg/database"
	"github.com/bdreece/herobrian/pkg/email"
	"github.com/bdreece/herobrian/pkg/identity"
	"github.com/bdreece/herobrian/pkg/linode"
	"github.com/bdreece/herobrian/pkg/systemd"
	"github.com/bdreece/herobrian/pkg/token"
	"github.com/bdreece/herobrian/web"
)

var (
	Config = fx.Module("config",
		fx.Provide(func(args Args) (config.Provider, error) {
			return config.NewYAML(
				config.Permissive(),
				config.Expand(os.LookupEnv),
				config.File(args.ConfigPath),
			)
		}),
	)

	Infrastructure = fx.Module("infrastructure",
		fx.Provide(
			logger.Configure,
			logger.New,
		),
		fx.Provide(
			database.Configure,
			database.ConfigureSuperUser,
			fx.Annotate(
				database.Dial,
				fx.As(new(database.DBTX)),
			),
			fx.Annotate(
				database.New,
				fx.As(new(database.Querier)),
				fx.As(fx.Self()),
			),
		),
		fx.Provide(
			email.ConfigureMailchimp,
			fx.Annotate(
				email.NewMailchimpClient,
				fx.As(new(email.Client)),
			),
		),
		// fx.Provide(
		// 	email.ConfigureSendGrid,
		// 	fx.Annotate(
		// 		email.NewSendGridClient,
		// 		fx.As(new(email.Client)),
		// 	),
		// ),
		fx.Provide(
			fx.Annotate(
				func(provider config.Provider) (*token.Options, error) {
					return token.Configure("password_reset", provider)
				},
				fx.ResultTags(`name:"password_reset"`),
			),
			fx.Annotate(
				token.NewHandler[token.PasswordResetClaims],
				fx.ParamTags(`name:"password_reset"`),
			),
			fx.Annotate(
				func(provider config.Provider) (*token.Options, error) {
					return token.Configure("user_invite", provider)
				},
				fx.ResultTags(`name:"user_invite"`),
			),
			fx.Annotate(
				token.NewHandler[token.UserInviteClaims],
				fx.ParamTags(`name:"user_invite"`),
			),
		),
		fx.Provide(
			identity.ConfigureSession,
			identity.ConfigureCookie,
			fx.Annotate(
				identity.NewSessionStore,
				fx.As(new(sessions.Store)),
			),
			fx.Annotate(
				identity.NewCookieAuthenticator,
				fx.As(new(identity.Authenticator)),
				fx.As(new(identity.SignInManager)),
			),
		),
		fx.Supply(
			fx.Annotate(
				echovalidator.Default,
				fx.As(new(echo.Validator)),
			),
		),
		fx.Provide(
			configureRenderer,
			fx.Annotate(
				echorenderer.New,
				fx.As(new(echo.Renderer)),
			),
		),
		fx.Provide(
			linode.Configure,
			linode.NewHTTP,
			linode.NewEmitter,
		),
		fx.Provide(
			systemd.ConfigureSSH,
			systemd.NewSSH,
			systemd.NewEmitter,
			systemd.NewServiceFactory,
		),
	)

	Application = fx.Module("application",
		fx.Provide(
			controller.NewHome,
			controller.NewAuth,
			controller.NewInvite,
			controller.NewLinode,
			controller.NewSystemd,
		),
		fx.Provide(
			router.Configure,
			router.New,
		),
		fx.Decorate(startRouter),
		fx.Decorate(createTables),
		fx.Invoke(func(router.Router) {}),
		fx.Invoke(func(*database.Queries) {}),
	)
)

func New(args Args) *fx.App {
	return fx.New(
		fx.Supply(args),
		Config,
		Infrastructure,
		Application,
	)
}

type Args struct {
	Port       int
	ConfigPath string
}

func (args Args) Addr() string {
	return fmt.Sprintf(":%d", args.Port)
}

func configureRenderer() *echorenderer.Options {
	return &echorenderer.Options{
		FS:      web.Templates,
		Include: []string{"*.gotmpl"},
		Funcs: func(c echo.Context) template.FuncMap {
			funcs := make(template.FuncMap)

			maps.Copy(funcs, sprig.FuncMap())
			maps.Copy(funcs, template.FuncMap{
				"claims": func() *identity.ClaimSet {
					claims, ok := c.Get(middleware.ClaimsContextKey).(*identity.ClaimSet)
					if !ok {
						return nil
					}

					return claims
				},
				"role": func() string {
					claims, ok := c.Get(middleware.ClaimsContextKey).(*identity.ClaimSet)
					if !ok {
						return ""
					}

					return identity.Role(claims.Role).String()
				},
			})

			return funcs
		},
	}
}

func startRouter(router router.Router, p struct {
	fx.In

	Home    *controller.Home
	Auth    *controller.Auth
	Invite  *controller.Invite
	Linode  *controller.Linode
	Systemd *controller.Systemd

	Args      Args
	Lifecycle fx.Lifecycle
}) router.Router {
	router.MapHome(p.Home)
	router.MapAuth(p.Auth)
	router.MapInvite(p.Invite)
	router.MapLinode(p.Linode)
	router.MapSystemd(p.Systemd)

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return router.Start(p.Args.Addr())
		},

		OnStop: func(ctx context.Context) error {
			return router.Shutdown(ctx)
		},
	})

	return router
}

func createTables(db *database.Queries, p struct {
	fx.In

	Options   *database.Options
	SuperUser *database.SuperUserOptions
	Lifecycle fx.Lifecycle
}) *database.Queries {
	p.Lifecycle.Append(fx.StartHook(func(ctx context.Context) (err error) {
		f, err := os.Open(p.Options.Schema)
		if err != nil {
			return
		}
		defer multierr.AppendInvoke(&err, multierr.Close(f))

		schema, err := io.ReadAll(f)
		if err != nil {
			return
		}

		_, err = db.DBTX().ExecContext(ctx, string(schema))
		if err != nil {
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(p.SuperUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return
		}

		_, err = db.UpsertUser(ctx, database.UpsertUserParams{
			ID:           -1,
			Username:     p.SuperUser.Username,
			PasswordHash: base64.StdEncoding.EncodeToString(hash),
			RoleID:       int64(identity.RoleSuper),
		})

		return
	}))

	return db
}
