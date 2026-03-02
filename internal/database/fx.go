package database

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/bdreece/herobrian/internal/identity"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	_ "modernc.org/sqlite"
)

var Module = fx.Module("database",
	fx.Provide(
		fx.Annotate(
			openDB,
			fx.As(fx.Self()),
			fx.As(new(DBTX)),
		),
		fx.Annotate(
			New,
			fx.As(fx.Self()),
			fx.As(new(Querier)),
		),
	),
	fx.Decorate(decorateQueries),
	fx.Invoke(func(*Queries) {}),
)

func openDB() (*sql.DB, error) {
	dsn := viper.GetString("sqlite:dsn")
	slog.Debug("connecting to database...")
	return sql.Open("sqlite", dsn)
}

func decorateQueries(queries *Queries, hasher identity.PasswordHasher, lc fx.Lifecycle) *Queries {
	lc.Append(fx.StartHook(func(ctx context.Context) error {
		slog.Debug("applying schema...")
		if _, err := queries.ApplySchema(ctx); err != nil {
			return err
		}

		var admin struct {
			FirstName   string `mapstructure:"first_name"`
			LastName    string `mapstructure:"last_name"`
			DisplayName string `mapstructure:"display_name"`
			Password    string `mapstructure:"password"`
			PictureURL  string `mapstructure:"picture_url"`
		}

		if err := viper.UnmarshalKey("sqlite:admin", &admin); err != nil {
			return err
		}

		digest, err := hasher.Hash(admin.Password)
		if err != nil {
			return err
		}

		params := UpsertUserParams{
			FirstName:    admin.FirstName,
			LastName:     admin.LastName,
			DisplayName:  admin.DisplayName,
			PictureURL:   &admin.PictureURL,
			PasswordHash: digest.String(),
		}

		slog.Debug("upserting admin", "params", params)

		if _, err := queries.UpsertUser(ctx, params); err != nil {
			slog.Error("failed to upsert admin", "error", err)
			return err
		}

		return nil
	}))

	return queries
}
