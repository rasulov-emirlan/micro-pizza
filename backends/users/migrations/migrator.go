package migrations

import (
	"database/sql"
	"embed"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// go:embed ./*.sql
var migrationFS embed.FS

func Up(url string) error {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}
	goose.SetBaseFS(migrationFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, ".")
}
