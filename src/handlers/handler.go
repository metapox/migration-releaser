package handlers

type DatabaseHandler interface {
	CreateDatabase(dbUrl string, dbName string) error
	UpMigrate(dbUrl string, path string) error
}

func NewDatabaseHandler(db string) DatabaseHandler {
	switch db {
	case "postgres":
		return &PostgresHandler{}
	case "mysql":
		return &MysqlHandler{}
	default:
		return nil
	}
}
