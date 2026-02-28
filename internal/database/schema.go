//go:generate go tool sqlc generate -f ../../configs/sqlc.yml
package database

import (
	"context"
	"database/sql"
	_ "embed"
)

//go:embed schema.sql
var schema string

func (q *Queries) ApplySchema(ctx context.Context) (sql.Result, error) {
	return q.db.ExecContext(ctx, schema)
}
