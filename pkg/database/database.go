//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate -f ../../configs/sqlc.yml
package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/config"
)

type SuperUserOptions struct {
	Username     string `yaml:"username"`
	FirstName    string `yaml:"first_name"`
	LastName     string `yaml:"last_name"`
	EmailAddress string `yaml:"email_address"`
	Password     string `yaml:"password"`
}

type Options struct {
	Schema           string `yaml:"schema"`
	ConnectionString string `yaml:"connection_string"`
}

func (q Queries) DBTX() DBTX { return q.db }

func ConfigureSuperUser(provider config.Provider) (*SuperUserOptions, error) {
	opts := new(SuperUserOptions)
	if err := provider.Get("database.super_user").Populate(opts); err != nil {
		return nil, fmt.Errorf("failed to configure database options: %w", err)
	}

	return opts, nil
}

func Configure(provider config.Provider) (*Options, error) {
	opts := new(Options)
	if err := provider.Get("database").Populate(opts); err != nil {
		return nil, fmt.Errorf("failed to configure database options: %w", err)
	}

	return opts, nil
}

func Dial(opts *Options) (*sql.DB, error) {
	return sql.Open("sqlite3", opts.ConnectionString)
}
