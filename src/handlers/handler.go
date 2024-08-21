package handlers

import "github.com/pkg/errors"

type DatabaseHandler interface {
	CreateDatabase(dbUrl string, dbName string) error
	UpMigrate(dbUrl string, path string) error
}

func NewDatabaseHandler(db string) (DatabaseHandler, error) {
	switch db {
	case "postgres":
		return &PostgresHandler{}, nil
	case "mysql":
		return &MysqlHandler{}, nil
	default:
		return nil, errors.WithStack(errors.New("unsupported database type " + db))
	}
}
