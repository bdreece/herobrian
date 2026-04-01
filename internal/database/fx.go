package database

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/bdreece/herobrian/internal/security"
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

func decorateQueries(queries *Queries, hasher security.PasswordHasher, lc fx.Lifecycle) *Queries {
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
			Role        string `mapstructure:"role"`
		}

		if err := viper.UnmarshalKey("sqlite:admin", &admin); err != nil {
			return err
		}

		digest, err := hasher.Hash(admin.Password)
		if err != nil {
			return err
		}

		userParams := UpsertUserParams{
			FirstName:    admin.FirstName,
			LastName:     admin.LastName,
			DisplayName:  admin.DisplayName,
			PictureURL:   &admin.PictureURL,
			PasswordHash: digest.String(),
		}

		if _, err := queries.UpsertUser(ctx, userParams); err != nil {
			slog.Error("failed to upsert admin", "error", err)
			return err
		}

		roleParams := []UpsertRoleParams{
			{Name: "landlubber"},
			{Name: "scallywag"},
			{Name: "freebooter"},
			{Name: "privateer"},
			{Name: "swashbuckler"},
		}

		for _, p := range roleParams {
			if _, err := queries.UpsertRole(ctx, p); err != nil {
				slog.Error("failed to upsert role", "error", err)
				return err
			}
		}

		return nil
	}))

	return queries
}
