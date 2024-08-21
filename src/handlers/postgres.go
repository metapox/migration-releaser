package handlers

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"log/slog"
)

type PostgresHandler struct {
}

func (p *PostgresHandler) CreateDatabase(dbUrl string, dbName string) error {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		if err.Error() != "pq: database \""+dbName+"\" already exists" {
			return errors.WithStack(err)
		} else {
			slog.Info("Database " + dbName + " already exists")
		}
	}

	return nil
}

func (p *PostgresHandler) UpMigrate(dbUrl string, s3path string) error {
	m, err := migrate.New(
		s3path,
		dbUrl)
	if err != nil {
		return errors.WithStack(err)
	}

	err = m.Up()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
