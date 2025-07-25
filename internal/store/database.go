package store

import (
	"GoBackendExploreMovieTracker/internal/utils"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
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
