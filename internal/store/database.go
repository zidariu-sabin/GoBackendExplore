package store

import (
	"GoBackendExploreMovieTracker/internal/utils"
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", utils.DATABASE_PGS_CONN_STRING)
	if err != nil {
		//format specifier that wraps error
		return nil, fmt.Errorf("db: open %w", err)
	}

	loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
	db = sqldblogger.OpenDriver(utils.DATABASE_PGS_CONN_STRING, db.Driver(), loggerAdapter /*, using_default_options*/)
	db.Ping()

	fmt.Println("Connected to Database... ")

	return db, nil
}

func MigrateFs(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)

}

// applying migrations to the database using Up function and checking for errors
func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)
	//we check if the database is up
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
