package handlers

type MysqlHandler struct {
}

func (m *MysqlHandler) CreateDatabase(dbUrl string, dbName string) error {
	return nil
}

func (m *MysqlHandler) UpMigrate(dbUrl string, path string) error {
	return nil
}
